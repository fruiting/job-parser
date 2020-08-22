<?php

namespace App\Jobs;

use App\Models\User;
use App\Models\Vacancy;
use App\Services\Parser\Parser;
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
class ParseWebSiteJob implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /** @var string $resource Web-site to parse */
    private $resource;

    /** @var Vacancy $vacancy Vacancy model */
    private $vacancy;

    /** @var User $user User model */
    private $user;

    /** @var int Max tries to run job */
    public $tries = 5;

    /**
     * Create a new job instance.
     *
     * @param string $resource Web-site to parse
     * @param Vacancy $vacancy Vacancy model
     * @param User $user User model
     */
    public function __construct(string $resource, Vacancy $vacancy, User $user)
    {
        $this->resource = $resource;
        $this->vacancy = $vacancy;
        $this->user = $user;
    }

    /**
     * Execute the job.
     *
     * @return void
     *
     * @throws \PHPHtmlParser\Exceptions\CircularException
     */
    public function handle()
    {
        (new Parser())->execute($this->resource, $this->vacancy, $this->user);
    }

    /**
     * Handles failed job
     *
     * @param null $exception
     */
    public function fail($exception = null)
    {
        logger()->error($exception->getMessage());
    }
}
