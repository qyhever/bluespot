export interface IMenuItem {
  title: string
  path: string
  icon?: string
  authKey?: string
  children?: IMenuItem[]
}

export const menuList: IMenuItem[] = [
  {
    title: '首页',
    path: '/home',
    icon: 'smile',
    // authKey: 'home',
  },
  {
    title: '组件',
    path: '/block',
    icon: 'smile',
    authKey: 'component',
    children: [
      {
        title: '复制',
        path: '/block/clipboard',
        icon: 'smile',
        authKey: 'clipboard',
      },
      {
        title: '二维码',
        path: '/block/qrcode',
        icon: 'smile',
        authKey: 'qrcode',
      },
    ],
  },
  {
    title: '上传',
    path: '/attach',
    icon: 'smile',
    authKey: 'component',
    children: [
      {
        title: '常规上传',
        path: '/attach/simple-upload',
        icon: 'smile'
      },
      {
        title: '分片上传',
        path: '/attach/chunk-upload',
        icon: 'smile'
      },
    ],
  },
  {
    title: '图表',
    path: '/chart',
    icon: 'smile',
    authKey: 'chart',
    children: [
      {
        title: '结婚率',
        path: '/chart/married',
        icon: 'smile',
        authKey: 'married',
      },
      {
        title: '出生率',
        path: '/chart/birth',
        icon: 'smile',
        authKey: 'birth',
      },
      {
        title: '高考',
        path: '/chart/college-entry-exam',
        icon: 'smile',
        authKey: 'college_entrance_examination',
      },
      {
        title: '大学毕业生',
        path: '/chart/university-graduate',
        icon: 'smile',
        authKey: 'university_graduate',
      },
      {
        title: '研究生报考',
        path: '/chart/postgraduate',
        icon: 'smile',
        authKey: 'postgraduate',
      },
    ],
  },
]
