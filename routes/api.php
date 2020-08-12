<?php

use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

Route::prefix('v3/parser')->group(function () {
    Route::post('/execute', 'ParserController@execute');
    Route::get('/overall/{userId}/{vacancyId}', 'ParserController@getOverall');
    Route::get('/vacancies/{userId}/{vacancyId}', 'ParserController@getVacancies');
});
