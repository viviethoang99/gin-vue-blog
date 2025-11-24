const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Message',
  path: '/message',
  component: Layout,
  redirect: '/message/comment',
  meta: {
    title: 'Message Management',
    icon: 'ic:twotone-email',
    order: 3,
    // role: ['admin'],
    // requireAuth: true,
  },
  children: [
    {
      name: 'CommentList',
      path: 'comment',
      component: () => import('./comment/index.vue'),
      meta: {
        title: 'Comment Management',
        icon: 'ic:twotone-comment',
        keepAlive: true,
      },
    },
    {
      name: 'LeaveMsgList',
      path: 'leave-msg',
      component: () => import('./leave-msg/index.vue'),
      meta: {
        title: 'Message Management',
        icon: 'ic:twotone-message',
        keepAlive: true,
      },
    },
  ],
}
