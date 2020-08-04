<?php

namespace App\Services\Parser;

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
     * @param string[] $skills Array of skills in vacancy
     *
     * @return void
     */
    public static function addSkills(array $skills): void
    {
        foreach ($skills as $skill) {
            $skillsCount = Redis::hget('romaspirin93@gmail.com:php-программист:skills', strtolower($skill));
            Redis::hset('romaspirin93@gmail.com:php-программист:skills', strtolower($skill), ++$skillsCount);
        }
    }

    /**
     * Returns array of popular skills
     *
     * @return int[]
     */
    public static function getOnlyPopular(): array
    {
        $skills = Redis::hgetall('romaspirin93@gmail.com:php-программист:skills');
        return array_filter($skills, function (string $skill) {
            return $skill >= self::COUNT_TO_BE_POPULAR;
        });
    }
}
