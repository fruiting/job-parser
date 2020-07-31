<?php

namespace App\Services\Parser;

use App\Services\Parser\HeadHunter\HeadHunterGeneralParser;

/**
 * Class ParserFactory describes logic of getting parser class
 *
 * @package App\Services\Parser
 */
class ParserFactory
{
    /**
     * Returns ParserInterface object by chosen web-site
     *
     * @param string $site
     *
     * @return GeneralParserInterface
     */
    public static function getParser(string $site): GeneralParserInterface
    {
        switch ($site) {
            default:
                $object = new HeadHunterGeneralParser();
                break;
        }

        return $object;
    }
}
