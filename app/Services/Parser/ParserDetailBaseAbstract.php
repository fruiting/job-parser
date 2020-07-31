<?php

namespace App\Services\Parser;

use PHPHtmlParser\Dom;
use PHPHtmlParser\Exceptions\ChildNotFoundException;
use PHPHtmlParser\Exceptions\CircularException;
use PHPHtmlParser\Exceptions\ContentLengthException;
use PHPHtmlParser\Exceptions\LogicalException;
use PHPHtmlParser\Exceptions\NotLoadedException;
use PHPHtmlParser\Exceptions\StrictException;
use Psr\Http\Client\ClientExceptionInterface;
use Throwable;

/**
 * Class ParserDetailBaseAbstract describes parse logic of detail page
 *
 * @package App\Services\Parser
 */
abstract class ParserDetailBaseAbstract implements DetailPageParserInterface
{
    /** @var Dom $dom Dom parser object */
    protected $dom;

    /** @var string $vacancyName */
    protected $vacancyName;

    /** @var string $companyName */
    protected $companyName;

    /** @var string $salaryText */
    protected $salaryText;

    /** @var float[] $salaryRange Array of salary range */
    protected $salaryRange;

    /** @var string[] $skills Array of skills in vacancy */
    protected $skills;

    /**
     * Parses vacancy detail page
     *
     * @param string $link Vacancy page link
     *
     * @return VacancyDto|null
     *
     * @throws ChildNotFoundException
     * @throws CircularException
     * @throws ContentLengthException
     * @throws LogicalException
     * @throws NotLoadedException
     * @throws StrictException
     * @throws ClientExceptionInterface
     */
    public function execute(string $link): ?VacancyDto
    {
        $this->dom = DomHelper::getInitedDom($link);

        $vacancy = null;
        try {
            $this->loadVacancyName();
            $this->loadCompany();
            $this->loadSalary();
            $this->loadSkills();
            $vacancy = new VacancyDto(
                $link,
                $this->vacancyName,
                $this->companyName,
                $this->salaryText,
                $this->salaryRange,
                $this->skills
            );
        } catch (Throwable
            | ChildNotFoundException
            | CircularException
            | ContentLengthException
            | LogicalException
            | NotLoadedException
            | StrictException
            | ClientExceptionInterface $exception) {
            //todo log it
        } finally {
            return $vacancy;
        }
    }
}
