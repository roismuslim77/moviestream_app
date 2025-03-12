<template>
  <section class="section">
    <div class="container">
      <h1 class="title">Movie List</h1>
    </div>
    <div class="filters row">
      <div class="form-group col-sm-3">
        <input
          v-model="searchKey"
          class="form-control"
          id="search-element"
          requred
        />
      </div>
    </div>
    <table class="table">
      <thead>
        <tr>
          <th></th>
          <th>Title</th>
          <th>Description</th>
          <th>Artist</th>
          <th>Genre</th>
          <th>Total Views | Votes</th>
          <th class="col-sm-2">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr></tr>
        <tr v-for="movie in movies" :key="movie.id">
          -->
          <td>
            <div>{{ movie.title }}</div>
          </td>
          <td>{{ movie.description }}</td>
          <td>{{ movie.artist }}</td>
          <td>{{ movie.genre }}</td>
          <td>{{ movie.total_views + " | " + movie.total_vote }}</td>
          <td></td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<script>
import Axios from "Axios";

export default {
  created(){
    this.fetchMovies()
  },
  layout: "vue-crud",
  data() {
    return { searchKey: "", movies: [] };
  },
  computed: {
    filteredtickets() {
      return this.movies.filter(
        movie =>
          movie.description
            .toLowerCase()
            .indexOf(this.searchKey.toLowerCase()) !== -1
      );
    }
  },
  methods: {
    async fetchMovies(){
        this.movies = await Axios.get('http://localhost:8000/movies', {
          headers: {
            Authorization: `${this.$auth.getToken('local') ? this.$auth.getToken('local').slice(7) : ''}`//'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl9pZCI6MSwidHlwZSI6ImN1c3RvbWVyIiwiZXhwIjoxNzQxODc0ODI0fQ.BrM6b0wab9R0mPqajBVEAKOML_r5u-mV3RHGILtVb98'
          },
        })
        .then((res) => res.data.data)
        .catch(() => []);

    },
    deleteTicketById(id) {
      let foundIndex = this.tickets.findIndex(p => p.id === id);
      Axios.delete("http://localhost:8080/ticket/?id=" + id).then(response => {
        this.tickets.splice(foundIndex, 1);
        console.log(this.tickets);
      });
    },
    updateTicket(ticket) {
      let foundIndex = this.tickets.findIndex(p => p.id === this.tickets.id);
      console.log("update " + foundIndex + this.tickets.name);
      //
    },
    sendMail(ticket) {
      let foundIndex = this.tickets.findIndex((p) => p.id === ticket.id);
      Axios.put("http://localhost:8080/sendMail/?email=" + ticket.user.email +"&status=" + ticket.status).then(
          (response) => {
            this.tickets.splice(foundIndex, 1);
            console.log(this.tickets);
          }
      );
    },
  },

  components: {}
};
</script>
<style>
.margin-5 {
  margin-right: 5px;
}
.title {
  text-align: center;
}
nav {
  display: flex;
  justify-content: flex-end;
}
.login {
  display: flex;
  align-items: center;

  justify-content: center;
}
.btn-control {
  display: flex;
}
.btn-edit {
  margin-right: 15px;
}

.form-group {
  max-width: 500px;
}

.actions {
  display: flex;
  justify-content: space-evenly;
  padding: 10px 0;
}

.glyphicon-euro {
  font-size: 12px;
}
</style>
