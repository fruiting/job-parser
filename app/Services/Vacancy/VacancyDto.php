<?php

namespace App\Services\Vacancy;

use Illuminate\Support\Fluent;

/**
 * Class VacancyDto describes vacancy entity
 *
 * @package App\Services\Parser
 *
 * @property-read string    $link           Vacancy detail page link
 * @property-read string    $vacancyName    Vacancy name
 * @property-read string    $companyName    Company name
 * @property-read string    $salaryText     Salary text (not formatted)
 * @property-read float[]   $salaryRange    Array of salary range of vacancy
 * @property-read string[]  $skills         Array of requirement skills of vacancy
 */
class VacancyDto extends Fluent
{
    /**
     * VacancyDto constructor.
     *
     * @param string $link Vacancy detail page link
     * @param string $vacancyName Vacancy name
     * @param string $companyName Company name
     * @param string $salaryText Salary text (not formatted)
     * @param float[] $salaryRange Array of salary range of vacancy
     * @param string[] $skills Requirement skills of vacancy
     */
    public function __construct(
        string $link,
        string $vacancyName,
        string $companyName,
        string $salaryText,
        array $salaryRange,
        array $skills
    ) {
        parent::__construct([
            'link'          => $link,
            'vacancyName'   => $vacancyName,
            'companyName'   => $companyName,
            'salaryText'    => $salaryText,
            'salaryRange'   => $salaryRange,
            'skills'        => $skills
        ]);
    }
}
