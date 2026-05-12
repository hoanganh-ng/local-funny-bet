import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth.store.js'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import LeaderboardView from '../views/LeaderboardView.vue'
import MatchView from '../views/MatchView.vue'
import JoinView from '../views/JoinView.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { public: true }
  },
  {
    path: '/',
    name: 'Home',
    component: HomeView,
    meta: { requiresAuth: true }
  },
  {
    path: '/matches/:id',
    name: 'Match',
    component: MatchView,
    meta: { requiresAuth: true }
  },
  {
    path: '/leaderboards',
    name: 'LeaderboardsList',
    component: () => import('../views/LeaderboardsListView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/leaderboards/new',
    name: 'CreateBoard',
    component: () => import('../views/CreateBoardView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/leaderboards/:id',
    name: 'Leaderboard',
    component: LeaderboardView,
    meta: { requiresAuth: true }
  },
  {
    path: '/history',
    name: 'History',
    component: () => import('../views/HistoryView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/join',
    name: 'Join',
    component: JoinView,
    meta: { public: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  if (to.path === '/auth/callback') {
    const success = to.query.success

    if (success === 'true') {
      try {
        // Exchange httpOnly refresh cookie for access token
        const response = await fetch('/api/auth/exchange', {
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' }
        })

        if (!response.ok) {
          next('/login?error=auth_failed')
          return
        }

        const data = await response.json()

        // Fetch user info
        const userResponse = await fetch('/api/auth/me', {
          headers: {
            'Authorization': `Bearer ${data.access_token}`,
            'Content-Type': 'application/json'
          }
        })

        if (userResponse.ok) {
          const user = await userResponse.json()
          authStore.setAuth(data.access_token, user)
        } else {
          authStore.setAuth(data.access_token, null)
        }

        next('/')
        return
      } catch (err) {
        console.error('OAuth callback error:', err)
        next('/login?error=auth_failed')
        return
      }
    } else {
      next('/login?error=auth_failed')
      return
    }
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router
