<?php

namespace App\Services\Parser;

/**
 * Interface GeneralParserInterface describes methods of general parser
 *
 * @package App\Services\Parser
 */
interface GeneralParserInterface
{
    /**
     * Returns object to parse list page
     *
     * @return ListPageParserInterface
     */
    public function getListPageParser(): ListPageParserInterface;

    /**
     * Returns object to parse detail page
     *
     * @return DetailPageParserInterface
     */
    public function getDetailPageParser(): DetailPageParserInterface;

    /**
     * Parses count of vacancies
     *
     * @param string $vacancyTitle
     *
     * @return int
     *
     * @throws \PHPHtmlParser\Exceptions\CircularException
     */
    public function getVacanciesCount (string $vacancyTitle): int;

    /**
     * Returns pages count
     *
     * @param string $vacanciesCount
     *
     * @return int
     */
    public function getPagesCount(string $vacanciesCount): int;
}
