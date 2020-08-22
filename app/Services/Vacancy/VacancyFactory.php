<?php

namespace App\Services\Vacancy;

/**
 * Class VacancyFactory
 *
 * @package App\Services\Vacancy
 */
class VacancyFactory
{
    /**
     * Returns VacancyDto builder from json string
     *
     * @param string $vacancyJson
     *
     * @return VacancyDto
     */
    public static function getFromJson(string $vacancyJson): VacancyDto
    {
        $vacancyObject = json_decode($vacancyJson);
        return new VacancyDto(
            $vacancyObject->link,
            $vacancyObject->vacancyName,
            $vacancyObject->companyName,
            $vacancyObject->salaryText,
            $vacancyObject->salaryRange,
            $vacancyObject->skills
        );
    }
}
