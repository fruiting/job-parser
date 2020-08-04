<?php

namespace App\Services\Parser;

use Generator;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Redis;

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
     * Executes parser
     *
     * @param string $site Site url
     * @param string $vacancyTitle Search by title
     * @param string $email User email
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\CircularException
     */
    public function execute(string $site, string $vacancyTitle, string $email): void
    {
        $factory = ParserFactory::getParser($site);
        $listPageParser = $factory->getListPageParser();
        $detailPageParser = $factory->getDetailPageParser();

        $vacanciesCount = $factory->getVacanciesCount($vacancyTitle);
        Redis::hset($email . ':' . $vacancyTitle, 'vacancies_count', $vacanciesCount);

        $pagesCount = $factory->getPagesCount($vacanciesCount);

        $generator = $this->parseDetails($pagesCount, $listPageParser, $detailPageParser, $vacancyTitle);
        $page = 1;
        foreach ($generator as $vacancies) {
            $i = 0;
            foreach ($vacancies as $vacancy) {
                /** @var VacancyDto $vacancy */

                Salary::addSalaries($vacancy->salaryRange);
                PopularSkills::addSkills($vacancy->skills);
                Redis::hset($email . ':' . $vacancyTitle . ':' . $page++, $i++, $vacancy->toJson());
            }die;
        }
    }
}
