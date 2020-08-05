<?php

namespace App\Http\Controllers;

use App\Jobs\ParseJobWebSite;
use App\Jobs\SendReportLink;
use App\Models\User;
use App\Models\Vacancy;
use App\Services\Parser\PopularSkills;
use App\Services\Parser\Salary;
use App\Services\Vacancy\VacancyRedis;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Response;
use Illuminate\Support\Collection;

/**
 * Class ParserController describes logic of parsing vacancies
 *
 * @package App\Http\Controllers
 */
class ParserController extends Controller
{
    /**
     * Executes parser
     *
     * @return JsonResponse
     *
     * @api
     */
    public function execute(): JsonResponse
    {
        $email = e(request()->get('email'));
        $vacancies = request()->get('vacancies');

        User::firstOrCreate(['email' => $email]);

        /** @var Collection|Vacancy[] $vacancies */
        $vacanciesCollection = Vacancy::whereIn('name', $vacancies)->get();

        foreach ($vacancies as $vacancy) {
            $vacancy = e($vacancy);

            dispatch(new ParseJobWebSite(request()->get('resource'), $vacancy));
            if ($vacanciesCollection->isEmpty())  {
                Vacancy::create(['name' => $vacancy]);
            }
        }

        dispatch(new SendReportLink($email));
        return response()->json([], Response::HTTP_OK);
    }

    /**
     * Returns overall vacancy info for report
     *
     * @param int $userId
     * @param int $vacancyId
     *
     * @return mixed[]
     */
    public function getOverall(int $userId, int $vacancyId): array
    {
        $key = VacancyRedis::getRedisKeyForVacation($userId, $vacancyId);
        $vacanciesCount = VacancyRedis::getVacanciesCount($key);
        $popularSkills = PopularSkills::getOnlyPopular($key);
        $salary = (new Salary())->loadSalary($key);

        return [
            'vacanciesCount' => $vacanciesCount,
            'popularSkills' => $popularSkills,
            'salaries' => [
                'minSalary' => $salary->getMinSalary(),
                'maxSalary' => $salary->getMaxSalary(),
                'averageSalary' => $salary->getAverageSalary()
            ]
        ];
    }

    /**
     * Returns array of vacancies jsons
     *
     * @param int $userId
     * @param int $vacancyId
     * @param int $page Page in redis
     *
     * @return string[]
     */
    public function getVacancies(int $userId, int $vacancyId, int $page): array
    {
        $key = VacancyRedis::getRedisKeyForVacation($userId, $vacancyId);
        return VacancyRedis::getFromHash($key . ':' . $page);
    }
}
