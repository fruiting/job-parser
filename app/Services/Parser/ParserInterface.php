<?php

namespace App\Services\Parser;

/**
 * Interface ParserInterface describes methods of parsing job web-sites
 *
 * @package App\Services\Parser
 */
interface ParserInterface
{
    /**
     * Executes parser
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\ContentLengthException
     * @throws \PHPHtmlParser\Exceptions\LogicalException
     * @throws \PHPHtmlParser\Exceptions\StrictException
     * @throws \Psr\Http\Client\ClientExceptionInterface
     */
    public function execute(): void;

    /**
     * Parses count of vacancies
     *
     * @return void
     */
    public function loadVacanciesCount(): void;

    /**
     * Parses all vacancies for description
     *
     * @return void
     */
    public function loadVacanciesInfo(): void;

    /**
     * Parses specific vacancy info
     *
     * @return void
     */
    public function loadSpecificVacancyInfo(): void;
}
