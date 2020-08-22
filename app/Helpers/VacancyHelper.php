<?php

namespace App\Helpers;

/**
 * Class VacancyHelper
 *
 * @package App\Helpers
 */
class VacancyHelper
{
    /** @var string Redis key postfix for all vacancies */
    public const VACANCIES_INFO_REDIS_KEY_POSTFIX = 'vacancies_info';

    /** @var string Redis key postfix for well paid vacancies */
    public const WELL_PAID_VACANCIES_INFO_REDIS_KEY_POSTFIX = 'well_paid_vacancies_info';
}
