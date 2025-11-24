const Layout = () => import('@/layout/index.vue')

export default {
  name: 'User',
  path: '/user',
  component: Layout,
  redirect: '/user/list',
  meta: {
    title: 'User Management',
    icon: 'ph:user-list-bold',
    order: 5,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'UserList',
      path: 'list',
      component: () => import('./list/index.vue'),
      meta: {
        title: 'User List',
        icon: 'mdi:account',
        keepAlive: true,
      },
    },
    {
      name: 'OnlineUserList',
      path: 'online',
      component: () => import('./online/index.vue'),
      meta: {
        title: 'Online Users',
        icon: 'ic:outline-online-prediction',
        keepAlive: true,
      },
    },
  ],
}
