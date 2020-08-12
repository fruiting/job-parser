<?php

namespace App\Console\Commands;

use App\Models\User;
use App\Models\Vacancy;
use App\Services\Parser\Parser as ParserObject;
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
    protected $signature = 'parser:run {vacancyTitle}';

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
        (new ParserObject())->execute('', Vacancy::where('name', $this->argument('vacancyTitle'))->first(), User::where('email', 'romaspirin93@gmail.com')->first());
        return 0;
    }
}
