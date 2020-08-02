<?php

namespace App\Mail;

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

    /**
     * SendReportLink constructor.
     *
     * @param string $link Report link
     */
    public function __construct(string $link)
    {
        $this->link = $link;
    }

    /**
     * Build the message.
     *
     * @return $this
     */
    public function build()
    {
        return $this//->from('no-reply@job-parser.ru')
            ->view('emails.report');
    }
}
