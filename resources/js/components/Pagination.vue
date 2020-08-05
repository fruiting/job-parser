<template>
    <nav class="mt-5 d-flex justify-content-center">
        <ul class="pagination">
            <li class="page-item">
                <a class="page-link" href="#" aria-label="Previous">
                    <span aria-hidden="true">&laquo;</span>
                    <span class="sr-only">Назад</span>
                </a>
            </li>
            <li class="page-item">
                <a class="page-link" href="javascript:;" @click="changePage(1)">
                    1
                </a>
            </li>
            <li v-if="this.pagesCount > 2" class="page-item">
                <a class="page-link" href="javascript:;">
                    ...
                </a>
            </li>
            <li class="page-item">
                <a class="page-link" href="javascript:;" @click="changePage(10)">
                    {{ this.pagesCount }}
                </a>
            </li>
            <li class="page-item">
                <a class="page-link" href="#" aria-label="Next">
                    <span aria-hidden="true">&raquo;</span>
                    <span class="sr-only">Вперед</span>
                </a>
            </li>
        </ul>
    </nav>
</template>

<script>
    export default {
        props: ['itemsCountOverall', 'itemsCountPerPage'],

        data() {
            return {
                pagesCount: 1
            };
        },

        mounted() {
            this.pagesCount = Math.ceil(this.$props.itemsCountOverall / this.$props.itemsCountPerPage);
        },

        methods: {
            changePage(page) {
            }
        }
    }
</script>
