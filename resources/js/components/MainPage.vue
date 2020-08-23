<template>
    <div class="container mt-5">
        <div v-if="parsingInProcess" class="d-flex align-items-center">
            <strong>
                Проверяю вакансии по запросам: {{ vacanciesTitlesImploded }}...<br>
                Эту страницу можно закрывать<br>
                На почту {{ email }} будут отправлены ссылка на отчеты после завершения парсинга.<br>
                Процесс займет некоторое время.
            </strong>
            <div class="spinner-border ml-auto" role="status" aria-hidden="true"></div>
        </div>
        <form v-else @submit="executeParser">
            <div class="form-group">
                <div class="row">
                    <div class="col">
                        <p>Вакансия</p>
                        <input v-model="vacanciesTitles[item - 1]" v-for="item in vacanciesCount" type="text"
                            class="form-control mb-2" required>
                    </div>
                    <div class="col">
                        <p>Куда направить ссылку на отчет</p>
                        <input type="email" v-model="email" class="form-control" aria-describedby="emailHelp"
                            placeholder="Email" required>
                        <button type="button" class="btn btn-primary mt-2" @click="addVacancy">
                            Добавить вакансию
                        </button>
                        <button type="submit" class="btn btn-success mt-2">
                            Запустить парсер
                        </button>
                    </div>
                </div>
            </div>
        </form>
    </div>
</template>

<script>
    export default {
        data() {
            return {
                /** @var {String[]} vacanciesTitles Array of vacancies titles to parse */
                vacanciesTitles: [],

                /** @var {String} vacanciesTitlesImploded Joined vacancies titles to parse */
                vacanciesTitlesImploded: null,

                /** @var {String} email Email to send link to report */
                email: null,

                /** @var {Integer} vacanciesCount Count of fields with vacancies titles */
                vacanciesCount: 1,

                /** @var {Integer} maxVacancies Max count of fields with vacancies titles */
                maxVacancies: 3,

                /** @var {Boolean} parsingInProcess Flag of working parser */
                parsingInProcess: false
            };
        },

        methods: {
            /**
             * Adds one more vacancy title field
             *
             * @return {VoidFunction}
             */
            addVacancy() {
                if (this.vacanciesCount !== this.maxVacancies) {
                    this.vacanciesCount++;
                }
            },

            /**
             * Sends request to backend to run parser
             *
             * @return {VoidFunction}
             */
            executeParser(e) {
                e.preventDefault();

                const data = {
                    resource: 'hh.ru',
                    vacancies: this.vacanciesTitles,
                    email: this.email
                };
                axios.post('/api/v3/parser/execute', data).then(() => {
                    this.vacanciesTitlesImploded = this.vacanciesTitles.join('; ');
                    this.parsingInProcess = true;
                });
            }
        }
    }
</script>
