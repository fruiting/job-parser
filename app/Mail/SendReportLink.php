<?php

namespace App\Mail;

use App\Models\Vacancy;
use Illuminate\Bus\Queueable;
use Illuminate\Mail\Mailable;
use Illuminate\Queue\SerializesModels;

/**
 * Class SendReportLink describes report mail settings
 *
 * @package App\Mail
 */
class SendReportLink extends Mailable
{
    use Queueable, SerializesModels;

    /** @var string $link Report link */
    public $link;

    /** @var Vacancy $vacancy */
    private $vacancy;

    /**
     * SendReportLink constructor.
     *
     * @param string $link Report link
     * @param Vacancy $vacancy
     */
    public function __construct(string $link, Vacancy $vacancy)
    {
        $this->link = $link;
        $this->vacancy = $vacancy;
    }

    /**
     * Build the message.
     *
     * @return $this
     */
    public function build()
    {
        return $this->from('no-reply@job-parser.ru')->view('emails.report');
    }
}
