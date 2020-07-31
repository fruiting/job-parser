<?php

namespace App\Services\Parser;

use Illuminate\Support\Collection;

/**
 * Class ParserBase
 *
 * @package App\Services\Parser
 */
final class Parser
{
    /**
     * Executes parser
     *
     * @param string $site Site url
     * @param string $vacancyTitle Search by title
     *
     * @return void
     */
    public function execute(string $site, string $vacancyTitle): void
    {
        $factory = ParserFactory::getParser($site);
        $listPageParser = $factory->getListPageParser();
        $detailPageParser = $factory->getDetailPageParser();

        $listPageParser->execute($vacancyTitle);
        $vacanciesUrls = $listPageParser->getVacanciesUrls();

        $vacanciesCollection = new Collection();
        foreach ($vacanciesUrls as $url) {
            /** @var VacancyDto $vacancy */
            $vacancy = $detailPageParser->execute($url);
            $vacanciesCollection->push($vacancy);
        }

        dd($vacanciesCollection);
    }
}
