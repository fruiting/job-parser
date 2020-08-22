<?php

namespace App\Services\Parser;

use App\Helpers\VacancyHelper;
use App\Models\User;
use App\Models\Vacancy;
use App\Services\Vacancy\VacancyDto;
use App\Services\Vacancy\VacancyFactory;
use App\Services\Vacancy\VacancyRedis;
use Generator;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Redis;
use Throwable;

/**
 * Class ParserBase
 *
 * @package App\Services\Parser
 */
final class Parser
{
    /**
     * Parses vacancies details
     *
     * @param int $pagesCount
     * @param ListPageParserInterface $listPageParser Parser object of list page
     * @param DetailPageParserInterface $detailPageParser Parser object of detail page
     * @param string $vacancyTitle
     *
     * @return Generator
     */
    private function parseDetails(
        int $pagesCount,
        ListPageParserInterface $listPageParser,
        DetailPageParserInterface $detailPageParser,
        string $vacancyTitle
    ): Generator {
        for ($i = 0; $i < $pagesCount; $i++) {
            $listPageParser->execute($vacancyTitle, $i);
            $vacanciesUrls = $listPageParser->getVacanciesUrls();
            $vacanciesCollection = new Collection();

            foreach ($vacanciesUrls as $url) {
                /** @var VacancyDto $vacancy */
                $vacancy = $detailPageParser->execute($url);
                $vacanciesCollection->push($vacancy);
            }

            yield $vacanciesCollection;
        }
    }

    /**
     * Loops all saved vacations to search the most well paid. Well paid - when salary is more or equal than average
     *
     * @param string $key Vacancies key
     *
     * @return void
     */
    private function loadWellPaidVacanciesInfo(string $key): void
    {
        $averageSalary = (new Salary())->loadSalary($key)->getAverageSalary();
        $vacancies = VacancyRedis::getFromHash($key . ':' . VacancyHelper::VACANCIES_INFO_REDIS_KEY_POSTFIX);

        $i = 0;
        foreach ($vacancies as $vacancy) {
            try {
                $vacancyDto = VacancyFactory::getFromJson($vacancy);
                foreach ($vacancyDto->salaryRange as $salary) {
                    if ($salary >= $averageSalary) {
                        VacancyRedis::saveVacancy(
                            $key . ':' . VacancyHelper::WELL_PAID_VACANCIES_INFO_REDIS_KEY_POSTFIX,
                            $i++,
                            $vacancyDto
                        );
                    }
                }
            } catch (Throwable $exception) {
                logger()
                    ->error('Could not write well paid vacancy to redis. Reason: ' . $exception->getMessage());
            }
        }
    }

    /**
     * Realizes parser logic
     *
     * @param string $site Site url
     * @param Vacancy $vacancy Vacancy model
     * @param User $user User model
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\CircularException
     */
    public function execute(string $site, Vacancy $vacancy, User $user): void
    {
        $key = $user->email . ':' . str_replace(' ', '+', $vacancy->name);
        $factory = ParserFactory::getParser($site);
        $listPageParser = $factory->getListPageParser();
        $detailPageParser = $factory->getDetailPageParser();

        $vacanciesCount = $factory->getVacanciesCount($vacancy->name);
        VacancyRedis::saveVacanciesCount($key, $vacanciesCount);

        $pagesCount = $factory->getPagesCount($vacanciesCount);
        $generator = $this->parseDetails($pagesCount, $listPageParser, $detailPageParser, $vacancy->name);

        $i = 0;
        foreach ($generator as $vacancies) {
            foreach ($vacancies as $vacancy) {
                try {
                    /** @var VacancyDto $vacancy */

                    Salary::addSalaries($key, $vacancy->salaryRange);
                    PopularSkills::addSkills($key, $vacancy->skills);
                    VacancyRedis::saveVacancy(
                        $key . ':' . VacancyHelper::VACANCIES_INFO_REDIS_KEY_POSTFIX,
                        $i++,
                        $vacancy
                    );
                } catch (Throwable $exception) {
                    logger()->error('Could not write vacancy to redis. Reason: ' . $exception->getMessage());
                }
            }
        }

        $this->loadWellPaidVacanciesInfo($key);
        Redis::del($key . ':' . VacancyHelper::VACANCIES_INFO_REDIS_KEY_POSTFIX, '*');
    }
}
