<?php

namespace App\Services\Parser;

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
     *
     * @return void
     */
    public function execute(string $site): void
    {
        $factory = ParserFactory::getParser($site);
        $listPageParser = $factory->getListPageParser();
        $detailPageParser = $factory->getDetailPageParser();

        $listPageParser->execute();
        $vacanciesUrls = $listPageParser->getVacanciesUrls();

        foreach ($vacanciesUrls as $url) {
            $detailPageParser->execute($url);
        }
    }
}
