<template>
  <div class="row">
    <div class="col-md-8">
      <div class="row">
        <div class="col-md-12">
          <h3>User followings</h3>
        </div>
      </div>
      <div class="row">
        <div v-if="users.length > 0" class="col-md-8">
          <div v-for="u in users" :key="u.id" :user="u">
            <UserCardList key="{{ u.id }}profileCard" :user="u" />
          </div>
        </div>
      </div>

      <div v-if="users.length > 0" class="row">
        <div class="col-md-12">
          <b-pagination-nav :link-gen="linkGen" :number-of-pages="totalPages"
                            use-router @change="onPageChanged"
          />
        </div>
      </div>
    </div>

    <div class="col-md-4">
      <UserCard key="{{ user.id}}profilePage" :user="user" />
      <Footer />
    </div>
  </div>
</template>

<script>
import get from 'lodash/get'
import UserCardList from '../user_card_list/user_card_list.vue'
import UserCard from '../../components/user_card/user_card.vue'
import Footer from '../../components/footer/footer.vue'

export default {
  components: {
    UserCardList,
    UserCard,
    Footer
  },
  data () {
    return {
      error: false,
      userId: null,
      users: [],
      usersError: '',
      usersLoaded: false,
      totalPages: 1
    }
  },
  computed: {
    user () {
      return this.$store.getters.findUser(this.userId)
    }
  },
  created () {
    const routeParams = this.$route.params
    this.currentPage = this.$route.query.page || 1
    this.load(routeParams.name || routeParams.id)
  },
  destroyed () {
  },
  methods: {
    load (userNameOrId) {
      console.debug('loading profile for ' + userNameOrId)
      const user = this.$store.getters.findUser(userNameOrId)
      if (user) {
        this.userId = user.id
        console.warn('load::user::nothing to do')
      } else {
        this.$store.dispatch('fetchUser', userNameOrId)
          .then(({ id }) => {
            this.userId = id
            console.warn('load::!user::fetchUser::id::nothing to do')
          })
          .catch((reason) => {
            console.warn('load::!user::fetchUser::!id')
            const errorMessage = get(reason, 'error.error')
            if (errorMessage) {
              this.error = errorMessage
            } else {
              this.error = 'Error loading user: ' + errorMessage
            }
          })
      }
      this.fetchRelationships()
    },
    fetchRelationships () {
      // Fetch followings
      const userId = this.userId
      this.$store.state.api.backendInteractor.fetchFriends({ id: userId, page: 1 })
        .then((users) => {
          this.users = users.items
          this.totalPages = users.totalPages
          this.currentPage = users.page
          this.usersLoaded = true
        })
        .catch((e) => {
          console.log('cannot fetch user followings: ' + e.message)
          this.usersError = e
          this.usersLoaded = false
        })
    },
    onPageChanged (page) {
      this.currentPage = page
      this.fetchRelationships()
    },
    linkGen (pageNum) {
      return pageNum === 1 ? '?' : `?page=${pageNum}`
    }
  }
}
</script>