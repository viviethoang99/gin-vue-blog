const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Article',
  path: '/article',
  component: Layout,
  redirect: '/article/list',
  meta: {
    title: 'Article Management',
    icon: 'ic:twotone-article',
    order: 2,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'ArticleList',
      path: 'list',
      component: () => import('./list/index.vue'),
      meta: {
        title: 'Article List',
        icon: 'material-symbols:format-list-bulleted',
        // role: ['admin'],
        // requireAuth: true,
        keepAlive: true,
      },
    },
    {
      name: 'ArticleWrite',
      path: 'write',
      component: () => import('./write/index.vue'),
      meta: {
        title: 'Publish Article',
        icon: 'icon-park-outline:write',
        // role: ['admin'],
        // requireAuth: true,
        keepAlive: true,
      },
    },
    {
      name: 'ArticleEdit',
      path: 'write/:id',
      component: () => import('./write/index.vue'),
      isHidden: true,
      meta: {
        title: 'Edit Article',
        icon: 'icon-park-outline:write',
        // role: ['admin'],
        // requireAuth: true,
        // keepAlive: true,
      },
    },
    {
      name: 'CategoryList',
      path: 'category-list',
      component: () => import('./category/index.vue'),
      meta: {
        title: 'Category Management',
        icon: 'tabler:category',
        // role: ['admin'],
        // requireAuth: true,
        keepAlive: true,
      },
    },
    {
      name: 'TagList',
      path: 'tag-list',
      component: () => import('./tag/index.vue'),
      meta: {
        title: 'Tag Management',
        icon: 'tabler:tag',
        keepAlive: true,
      },
    },
  ],
}
