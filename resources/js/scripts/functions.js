const functions = {
    /**
     * Returns value of $_GET parameter from string
     *
     * @param string
     * @param parameter
     *
     * @return {string|null}
     */
    getParameterFromString(string, parameter) {
        const resultArray = new RegExp('[\?&]' + parameter + '=([^&#]*)').exec(string);
        let result = null;
        if (resultArray) {
            result = resultArray[1] || 0;
        }

        return result;
    }
};

export default functions;
