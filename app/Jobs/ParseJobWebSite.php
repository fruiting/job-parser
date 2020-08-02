<?php

namespace App\Jobs;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;

/**
 * Class ParseJobWebSite describes job for parsing job web-site
 *
 * @package App\Jobs
 */
class ParseJobWebSite implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /** @var string $resource Web-site to parse */
    private $resource;

    /** @var string $vacancy Vacancy name */
    private $vacancy;

    /** @var int Max tries to run job */
    public $tries = 5;

    /**
     * Create a new job instance.
     *
     * @param string $resource Web-site to parse
     * @param string $vacancy Vacancy name
     */
    public function __construct(string $resource, string $vacancy)
    {
        $this->resource = $resource;
        $this->vacancy = $vacancy;
    }

    /**
     * Execute the job.
     *
     * @return void
     */
    public function handle()
    {
        //
    }
}
