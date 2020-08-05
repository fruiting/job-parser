<?php

namespace App\Services\Parser;

use Illuminate\Support\Facades\Redis;

/**
 * Class Salary describes salaries in vacancies
 *
 * @package App\Services\Parser
 */
class Salary
{
    /** @var float[] $salaries Array of salaries */
    private $salaries;

    /**
     * Writes salaries in redis
     *
     * @param string $key Redis key
     * @param float[] $salaries Array of salaries in vacancy
     *
     * @return void
     */
    public static function addSalaries(string $key, array $salaries): void
    {
        foreach ($salaries as $salary) {
            Redis::lpush($key . ':salaries', $salary);
        }
    }

    /**
     * Loads all salaries of vacancy into object
     *
     * @param string $key Redis key
     *
     * @return $this
     */
    public function loadSalary(string $key): self
    {
        $this->salaries = Redis::lrange($key . ':salaries', 0, -1);
        return $this;
    }

    /**
     * Returns minimal salary of vacancy
     *
     * @return float
     */
    public function getMinSalary(): float
    {
        return min($this->salaries);
    }

    /**
     * Returns maximum salary of vacancy
     *
     * @return float
     */
    public function getMaxSalary(): float
    {
        return max($this->salaries);
    }

    /**
     * Returns average salary of vacancy
     *
     * @return float
     */
    public function getAverageSalary(): float
    {
         $values = array_count_values($this->salaries);
         arsort($values);
         return array_slice(array_keys($values), 0, 1, true)[0];
    }
}
