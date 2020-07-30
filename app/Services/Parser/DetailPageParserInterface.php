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
     * Parses specific vacancy info
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\NotLoadedException
     */
    public function loadVacancyInfo(): void;
}
