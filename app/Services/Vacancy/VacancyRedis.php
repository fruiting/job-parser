<?php

namespace App\Services\Vacancy;

use App\Models\User;
use App\Models\Vacancy;
use Illuminate\Support\Facades\Redis;

/**
 * Class VacancyRedis
 *
 * @package App\Services\Vacancy
 */
class VacancyRedis
{
    /** @var string Redis vacancies count object */
    private const VACANCIES_COUNT_OBJECT = 'vacancies_count';

    /**
     * Returns redis key for vacancy info
     *
     * @param User $user User model
     * @param Vacancy $vacancy Vacancy model
     *
     * @return string
     */
    public static function getRedisKeyForVacation(User $user, Vacancy $vacancy): string
    {
        return $user->email . ':' . $vacancy->name;
    }

    /**
     * Returns from redis hash table
     *
     * @param string $key Redis key for hash
     *
     * @return string[]
     */
    public static function getFromHash(string $key): array
    {
        return Redis::hgetall($key);
    }

    /**
     * Saves vacancies count to redis
     *
     * @param string $key Redis key
     * @param int $vacanciesCount
     *
     * @return void
     */
    public static function saveVacanciesCount(string $key, int $vacanciesCount): void
    {
        Redis::hset($key, self::VACANCIES_COUNT_OBJECT, $vacanciesCount);
    }

    /**
     * Returns vacancies count from redis
     *
     * @param string $key Redis key
     *
     * @return int
     */
    public static function getVacanciesCount(string $key): int
    {
        return Redis::hget($key, self::VACANCIES_COUNT_OBJECT);
    }

    /**
     * Saves vacancy object to redis
     *
     * @param string $key Redis key
     * @param int $i Object
     * @param VacancyDto $vacancy Vacancy object
     *
     * @return void
     */
    public static function saveVacancy(string $key, int $i, VacancyDto $vacancy): void
    {
        Redis::hset($key, $i, $vacancy->toJson());
    }
}
