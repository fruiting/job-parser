<?php

namespace App\Services\Parser;

/**
 * Class LinkedInParser describes logic of parsing ru.linkedin.com
 *
 * @package App\Services\Parser
 */
class LinkedInParser extends ParserBaseAbstract
{
    /** @var string Web-site link to parse */
    public const LINK = 'ru.linkedin.com';

    /**
     * Parses count of vacancies
     *
     * @return void
     */
    public function loadVacanciesCount(): void
    {
        // TODO: Implement loadVacanciesCount() method.
    }

    /**
     * Parses all vacancies for description
     *
     * @return void
     */
    public function loadVacanciesInfo(): void
    {
        // TODO: Implement loadVacanciesInfo() method.
    }

    /**
     * Parses specific vacancy info
     *
     * @return void
     */
    public function loadSpecificVacancyInfo(): void
    {
        // TODO: Implement loadSpecificVacancyInfo() method.
    }
}