<?php

namespace Tests\Unit;

use App\Models\Vacancy;
use App\Services\Parser\ParserFactory;
use Tests\TestCase;

/**
 * Class JobWebSiteConnectionTest
 *
 * @package Tests\Unit
 */
class JobWebSiteConnectionTest extends TestCase
{
    /**
     * Tests that web site has some vacations
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\CircularException
     */
    public function testApplication(): void
    {
        $vacancy = Vacancy::first();
        $factory = ParserFactory::getParser('hh.ru');
        $vacanciesCount = $factory->getVacanciesCount($vacancy->name);

        $this->assertGreaterThan(0, $vacanciesCount);
    }
}
