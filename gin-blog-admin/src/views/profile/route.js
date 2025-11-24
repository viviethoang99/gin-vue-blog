const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Profile',
  path: '/',
  component: Layout,
  redirect: '/profile',
  meta: {
    order: 8,
  },
  isCatalogue: true,
  children: [
    {
      name: 'Profile',
      path: '/profile',
      component: () => import('./index.vue'),
      meta: {
        title: 'Profile',
        icon: 'mdi:account',
        order: 0,
      },
    },
  ],
}
