<?php

namespace App\Services\Parser;

/**
 * Interface DetailPageParserInterface describes methods of parsing vacancies detail pages of job web-sites
 *
 * @package App\Services\Parser
 */
interface DetailPageParserInterface
{
    /**
     * Loads vacancy name
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadVacancyName(): void;

    /**
     * Loads salary info
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadSalary(): void;

    /**
     * Loads company name
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadCompany(): void;

    /**
     * Loads requirement skills of vacancy
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadSkills(): void;
}
