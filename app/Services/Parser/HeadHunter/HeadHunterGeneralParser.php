<?php

namespace App\Services\Parser\HeadHunter;

use App\Services\Parser\DetailPageParserInterface;
use App\Services\Parser\DomHelper;
use App\Services\Parser\ListPageParserInterface;
use App\Services\Parser\ParserGeneralBase;
use PHPHtmlParser\Dom\Node\HtmlNode;

/**
 * Class HeadHunterGeneralParser describes general
 *
 * @package App\Services\Parser\HeadHunter
 */
class HeadHunterGeneralParser extends ParserGeneralBase
{
    /** @var int Vacancies count on page */
    protected const PAGE_SIZE = 50;

    /**
     * Returns object to parse list page
     *
     * @return ListPageParserInterface
     */
    public function getListPageParser(): ListPageParserInterface
    {
        return new HeadHunterListPageParser();
    }

    /**
     * Returns object to parse detail page
     *
     * @return DetailPageParserInterface
     */
    public function getDetailPageParser(): DetailPageParserInterface
    {
        return new HeadHunterDetailPageParser();
    }

    /**
     * Parses count of vacancies
     *
     * @param string $vacancyTitle
     *
     * @return int
     *
     * @throws \PHPHtmlParser\Exceptions\CircularException
     */
    public function getVacanciesCount(string $vacancyTitle): int
    {
        $dom = DomHelper::getInitedDom(HeadHunterListPageParser::LINK . $vacancyTitle);

        /** @var HtmlNode $html */
        $html = $dom->find('h1');
        $header = $html->getChildren()[0];
        preg_match('!\d+!', $header->text(), $matches);
        return (int) $matches[0];
    }
}
