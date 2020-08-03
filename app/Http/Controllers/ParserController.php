<?php

namespace App\Http\Controllers;

use App\Jobs\ParseJobWebSite;
use App\Jobs\SendReportLink;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Response;
use Illuminate\Support\Facades\Redis;

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
        $vacancies = request()->get('vacancies');
//        foreach ($vacancies as $vacancy) {
//            dispatch(new ParseJobWebSite(request()->get('resource'), $vacancy));
//        }
        dispatch(new SendReportLink(e(request()->get('email'))));
        return response()->json([], Response::HTTP_OK);
    }

    public function get()
    {
        return Redis::hgetall('romaspirin93@gmail.com:php-программист');
    }
}
