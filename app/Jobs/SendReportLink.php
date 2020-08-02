<?php

namespace App\Jobs;

use App\Mail\SendReportLink as SendReportLinkMail;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Mail;

/**
 * Class SendReportLink describes job that sends report to email
 *
 * @package App\Jobs
 */
class SendReportLink implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /** @var string $email Email to send report link */
    private $email;

    /**
     * Create a new job instance.
     *
     * @param string $email Email to send report link
     */
    public function __construct(string $email)
    {
        $this->email = $email;
    }

    /**
     * Execute the job.
     *
     * @return void
     */
    public function handle()
    {
        //todo сделать получение ссылки
        Mail::to($this->email)->send(new SendReportLinkMail('https://google.ru'));
    }
}
