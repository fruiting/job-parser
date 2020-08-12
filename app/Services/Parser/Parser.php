<?php

namespace App\Services\Parser;

use App\Models\User;
use App\Models\Vacancy;
use App\Services\Vacancy\VacancyDto;
use App\Services\Vacancy\VacancyRedis;
use Generator;
use Illuminate\Support\Collection;

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
        for ($i = 0; $i < 2; $i++) {
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
     * Executes parser
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
        $key = $user->email . ':' . $vacancy->name;
        $factory = ParserFactory::getParser($site);
        $listPageParser = $factory->getListPageParser();
        $detailPageParser = $factory->getDetailPageParser();

        $vacanciesCount = $factory->getVacanciesCount($vacancy->name);
        VacancyRedis::saveVacanciesCount($key, $vacanciesCount);

        $pagesCount = $factory->getPagesCount($vacanciesCount);
        $generator = $this->parseDetails($pagesCount, $listPageParser, $detailPageParser, $vacancy->name);
        $page = 1;

        foreach ($generator as $vacancies) {
            $i = 0;
            foreach ($vacancies as $vacancy) {
                /** @var VacancyDto $vacancy */

                Salary::addSalaries($key, $vacancy->salaryRange);
                PopularSkills::addSkills($key, $vacancy->skills);
                VacancyRedis::saveVacancy($key . ':' . $page, $i, $vacancy);

                $page++;
                $i++;
            }
            break;
        }
    }
}
