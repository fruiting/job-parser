<template>
    <div class="container mt-5 mb-5">
        <h1>{{ vacancyName }}</h1>
        <div class="row mt-3">
            <div class="col-md-4">
                <h3>Количество вакансий</h3>
                <p>{{ vacanciesCount }}</p>
            </div>
            <div class="col-md-4">
                <h3>З/п</h3>
                Ср: {{ averageSalary }}<br>
                Мин: {{ minSalary }}<br>
                Макс: {{ maxSalary }}
            </div>
            <div class="col-md-4">
                <a :href="vacanciesListLink" target="_blank">Перейти к списку вакансий</a>
            </div>
        </div>

        <div class="mt-5">
            <h3>Наиболее частые скиллы</h3>
            <ul class="list-group">
                <li v-for="(count, skill) in popularSkills"
                    class="list-group-item d-flex justify-content-between align-items-center">
                    {{ skill }}
                    <span class="badge badge-primary badge-pill">
                        {{ count }}
                    </span>
                </li>
            </ul>
        </div>

        <div class="mt-5">
            <h3>Самые оплачиваемые вакансии</h3>
            <vacancy v-for="vacancy in vacancies" :vacancy="vacancy"></vacancy>
        </div>
    </div>
</template>

<script>
    import functions from '../scripts/functions';

    export default {
        data() {
            return {
                /** @var {String} vacanciesListLink */
                vacanciesListLink: null,

                /** @var {Integer} userId */
                userId: null,

                /** @var {Integer} vacancyId */
                vacancyId: null,

                /** @var {String} vacancyName */
                vacancyName: null,

                /** @var {Integer} vacanciesCount */
                vacanciesCount: 0,

                /** @var {Float} minSalary */
                minSalary: 0,

                /** @var {Float} maxSalary */
                maxSalary: 0,

                /** @var {Float} averageSalary */
                averageSalary: 0,

                /** @var {Mixed[]} vacancies Array of vacancies */
                vacancies: [],

                /** @var {String[]} popularSkills Array of the most often required skills in vacancies */
                popularSkills: []
            };
        },

        mounted() {
            this.userId = functions.getParameterFromString(window.location.search, 'userId');
            this.vacancyId = functions.getParameterFromString(window.location.search, 'vacancyId');

            axios.get('/api/v3/parser/overall/' + this.userId + '/' + this.vacancyId).then((response) => {
                this.vacancyName = response.data.vacancyName;
                this.vacanciesListLink = response.data.vacanciesListLink;
                this.vacanciesCount = response.data.vacanciesCount;
                this.popularSkills = response.data.popularSkills;
                this.minSalary = response.data.salaries.minSalary;
                this.maxSalary = response.data.salaries.maxSalary;
                this.averageSalary = response.data.salaries.averageSalary;
            });

            this.renderVacanciesList();
        },

        methods: {
            /**
             * Makes request to backend to load vacancies and updates component model
             *
             * @return {VoidFunction}
             */
            renderVacanciesList() {
                axios.get('/api/v3/parser/vacancies/' + this.userId + '/' + this.vacancyId)
                    .then((response) => {
                        for (let index in response.data) {
                            this.vacancies.push(JSON.parse(response.data[index]));
                        }
                    });
            }
        }
    }
</script>
