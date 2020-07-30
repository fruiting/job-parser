<?php

namespace App\Services\Parser;

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
     * @return ParserInterface
     */
    public static function getParser(string $site): ParserInterface
    {
        switch ($site) {
            case LinkedInParser::LINK:
                $object = new LinkedInParser();
                break;
            default:
                $object = new HeadHunterParser();
                break;
        }

        return $object;
    }
}
