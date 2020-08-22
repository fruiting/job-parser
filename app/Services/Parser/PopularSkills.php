<?php

namespace App\Services\Parser;

use App\Services\Vacancy\VacancyRedis;
use Illuminate\Support\Facades\Redis;

/**
 * Class PopularSkills describes skills in vacancies
 *
 * @package App\Services\Parser
 */
class PopularSkills
{
    /** @var int Minimum count to be popular skill */
    private const COUNT_TO_BE_POPULAR = 10;

    /**
     * Writes in redis vacancies skills
     *
     * @param string $key Redis key
     * @param string[] $skills Array of skills in vacancy
     *
     * @return void
     */
    public static function addSkills(string $key, array $skills): void
    {
        foreach ($skills as $skill) {
            $skillsCount = Redis::hget($key . ':skills', strtolower($skill));
            Redis::hset($key . ':skills', strtolower($skill), ++$skillsCount);
        }
    }

    /**
     * Returns array of popular skills
     *
     * @param string $key Redis key
     *
     * @return int[]
     */
    public static function getOnlyPopular(string $key): array
    {
        $skills = VacancyRedis::getFromHash($key . ':skills');
        return array_filter($skills, function (string $skill) {
            return $skill >= self::COUNT_TO_BE_POPULAR;
        });
    }
}
