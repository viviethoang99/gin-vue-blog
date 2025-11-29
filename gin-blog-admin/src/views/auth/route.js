const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Auth',
  path: '/auth',
  component: Layout,
  redirect: '/auth/menu',
  meta: {
    title: 'Permission Management',
    icon: 'cib:adguard',
    order: 3,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'MenuList',
      path: 'menu',
      component: () => import('./menu/index.vue'),
      meta: {
        title: 'Menu Management',
        icon: 'ic:twotone-menu-book',
        keepAlive: true,
      },
    },
    {
      name: 'ResourceList',
      path: 'resource',
      component: () => import('./resource/index.vue'),
      meta: {
        title: 'API Management',
        icon: 'mdi:api',
        keepAlive: true,
      },
    },
    {
      name: 'RoleList',
      path: 'role',
      component: () => import('./role/index.vue'),
      meta: {
        title: 'Role Management',
        icon: 'carbon:user-role',
        keepAlive: true,
      },
    },
  ],
}
