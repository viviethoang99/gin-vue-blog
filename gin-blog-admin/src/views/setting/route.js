const Layout = () => import('@/layout/index.vue')

export default {
  name: 'System',
  path: '/system',
  component: Layout,
  redirect: '/system/website',
  meta: {
    title: 'System Management',
    icon: 'ion:md-settings',
    order: 6,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'Website',
      path: 'website',
      component: () => import('./website/index.vue'),
      meta: {
        title: 'Website Management',
        icon: 'el:website',
        order: 1,
        keepAlive: true,
      },
    },
    {
      name: 'Page Management',
      path: 'page',
      component: () => import('./page/index.vue'),
      meta: {
        title: 'Page Management',
        icon: 'iconoir:journal-page',
        order: 2,
        keepAlive: true,
      },
    },
    {
      name: 'FriendLink',
      path: 'link',
      component: () => import('./link/index.vue'),
      meta: {
        title: 'Friend Links',
        icon: 'mdi:telegram',
        order: 3,
        keepAlive: true,
      },
    },
    {
      name: 'About',
      path: 'about',
      component: () => import('./about/index.vue'),
      meta: {
        title: 'About Me',
        icon: 'cib:about-me',
        order: 4,
        keepAlive: true,
      },
    },
  ],
}
