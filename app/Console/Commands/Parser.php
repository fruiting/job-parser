<?php

namespace App\Console\Commands;

use App\Services\Parser\ParserFactory;
use Illuminate\Console\Command;

/**
 * Class Parser which parses specific jobs web-site
 *
 * @package App\Console\Commands
 */
class Parser extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'parser:run';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Parses job web-sites for vacancies';

    /**
     * Execute the console command.
     *
     * @return int
     *
     * @throws \PHPHtmlParser\Exceptions\ChildNotFoundException
     * @throws \PHPHtmlParser\Exceptions\CircularException
     * @throws \PHPHtmlParser\Exceptions\ContentLengthException
     * @throws \PHPHtmlParser\Exceptions\LogicalException
     * @throws \PHPHtmlParser\Exceptions\StrictException
     * @throws \Psr\Http\Client\ClientExceptionInterface
     */
    public function handle()
    {
        ParserFactory::getParser('https://hh.ru/search/vacancy?st=searchVacancy&text=PHP+%D0%BF%D1%80%D0%BE%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B8%D1%81%D1%82&search_field=name&area=1&salary=&currency_code=RUR&experience=doesNotMatter&order_by=relevance&search_period=&items_on_page=50&no_magic=true&L_save_area=true&from=suggest_post')->execute();
        return 0;
    }
}
