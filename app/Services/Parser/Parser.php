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
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\CircularException
     */
    public function execute(string $site, string $vacancyTitle): void
    {
//        dd(Redis::hgetall('romaspirin93@gmail.com:' . $vacancyTitle));

        $factory = ParserFactory::getParser($site);
        $listPageParser = $factory->getListPageParser();
        $detailPageParser = $factory->getDetailPageParser();

        $vacanciesCount = $factory->getVacanciesCount($vacancyTitle);
        $pagesCount = $factory->getPagesCount($vacanciesCount);

        $generator = $this->parseDetails($pagesCount, $listPageParser, $detailPageParser, $vacancyTitle);
        foreach ($generator as $vacancies) {
            $i = 0;
            foreach ($vacancies as $vacancy) {
                /** @var VacancyDto $vacancy */
                Redis::hset('romaspirin93@gmail.com:' . $vacancyTitle, $i++, $vacancy->toJson());
            }
            die;
        }
    }
}
