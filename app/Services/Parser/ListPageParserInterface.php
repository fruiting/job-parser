<?php

namespace App\Services\Parser;

use PHPHtmlParser\Exceptions\ChildNotFoundException;
use PHPHtmlParser\Exceptions\NotLoadedException;

/**
 * Interface ListPageParserInterface describes methods of parsing vacancies list page of job web-sites
 *
 * @package App\Services\Parser
 */
interface ListPageParserInterface
{
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
     *
     * @throws ChildNotFoundException
     * @throws NotLoadedException
     */
    public function loadVacanciesInfo(): void;

    /**
     * Returns count of vacancies
     *
     * @return int
     */
    public function getVacanciesCount(): int;

    /**
     * Returns array of urls of vacancies
     *
     * @return string[]
     */
    public function getVacanciesUrls(): array;
}
