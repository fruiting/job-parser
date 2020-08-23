<?php

namespace App\Http\Controllers;

use App\Helpers\VacancyHelper;
use App\Jobs\ParseWebSiteJob;
use App\Models\User;
use App\Models\Vacancy;
use App\Services\Parser\HeadHunter\HeadHunterListPageParser;
use App\Services\Parser\PopularSkills;
use App\Services\Parser\Salary;
use App\Services\Vacancy\VacancyRedis;
use App\Services\Vacancy\VacancyUserRepository;
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
     */
    public function execute(): JsonResponse
    {
        $email = e(request()->get('email'));
        $vacancies = request()->get('vacancies');

        $user = User::firstOrCreate(['email' => $email]);

        /** @var Collection|Vacancy[] $vacancies */
        $vacanciesCollection = Vacancy::whereIn('name', $vacancies)->get();

        foreach ($vacancies as $key => $vacancyName) {
            $vacancyName = e($vacancyName);

            $vacancyModel = $vacanciesCollection->first(function (Vacancy $vacancy) use ($vacancyName) {
                return $vacancy->name == $vacancyName;
            });
            if (!$vacancyModel) {
                $vacancyModel = Vacancy::create(['name' => $vacancyName]);
            }

            ParseWebSiteJob::dispatch(request()->get('resource'), $vacancyModel, $user)->onQueue('vacancy-' . $key);
        }

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
        $repository = (new VacancyUserRepository())->loadData($userId, $vacancyId);
        $key = VacancyRedis::getRedisKeyForVacation($repository->getUser(), $repository->getVacancy());
        $vacanciesCount = VacancyRedis::getVacanciesCount($key);
        $popularSkills = PopularSkills::getOnlyPopular($key);
        $salary = (new Salary())->loadSalary($key);

        return [
            'vacancyName' => $repository->getVacancy()->name,
            'vacanciesListLink' => HeadHunterListPageParser::LINK . $repository->getVacancy()->name,
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
     * Returns array of well paid vacancies jsons
     *
     * @param int $userId
     * @param int $vacancyId
     *
     * @return string[]
     */
    public function getVacancies(int $userId, int $vacancyId): array
    {
        $repository = (new VacancyUserRepository())->loadData($userId, $vacancyId);
        $key = VacancyRedis::getRedisKeyForVacation($repository->getUser(), $repository->getVacancy());
        return VacancyRedis::getFromHash($key . ':' . VacancyHelper::WELL_PAID_VACANCIES_INFO_REDIS_KEY_POSTFIX);
    }
}
