<?php

namespace App\Services\Vacancy;

use App\Models\User;
use App\Models\Vacancy;

/**
 * Class VacancyUserRepository
 *
 * @package App\Services\Vacancy
 */
class VacancyUserRepository
{
    /** @var User $user */
    private $user;

    /** @var Vacancy $vacancy */
    private $vacancy;

    /**
     * Loads user and vacancy models to object
     *
     * @param int $userId User id in database
     * @param int $vacancyId Vacancy id in database
     *
     * @return self
     */
    public function loadData(int $userId, int $vacancyId): self
    {
        $this->user = User::where('id', $userId)->first();
        $this->vacancy = Vacancy::where('id', $vacancyId)->first();

        return $this;
    }

    /**
     * Returns user model
     *
     * @return User
     */
    public function getUser(): User
    {
        return $this->user;
    }

    /**
     * Returns vacancy model
     *
     * @return Vacancy
     */
    public function getVacancy(): Vacancy
    {
        return $this->vacancy;
    }
}
