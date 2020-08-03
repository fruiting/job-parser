<template>
    <div class="container mt-5">
    </div>
</template>

<script>
    export default {
        data() {
            return {
                /** @var {Integer} vacanciesCount */
                vacanciesCount: 0,

                /** @var {Float} minSalary */
                minSalary: 0,

                /** @var {Float} maxSalary */
                maxSalary: 0,

                /** @var {Mixed[]} vacancies Array of vacancies */
                vacancies: [],
            };
        },

        mounted() {
            axios.get('/api/v3/parser').then((response) => {
                for (let index in response.data) {
                    this.vacancies.push(JSON.parse(response.data[index]));
                }
            });
        }
    }
</script>
