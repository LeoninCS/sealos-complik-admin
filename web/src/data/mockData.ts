import type {
  ActivityItem,
  BanRecord,
  CommitmentRecord,
  ConfigRecord,
  NamespaceProfile,
  QuickLinkItem,
  StatCardItem,
  UnbanRecord,
  ViolationRecord,
} from "../types";

export const stats: StatCardItem[] = [
  {
    label: "违规 namespace 数",
    value: "24",
    delta: "24 条违规记录",
    tone: "danger",
    description: "按当前违规记录对应的 namespace 去重统计",
    targetPath: "/violations",
  },
  {
    label: "当前封禁数",
    value: "7",
    delta: "+1",
    tone: "warn",
    description: "当前仍处于有效封禁状态的 namespace",
    targetPath: "/bans",
  },
  {
    label: "今日新增违规",
    value: "12",
    delta: "+3",
    tone: "info",
    description: "两类违规总数，按事件统计",
    targetPath: "/violations",
  },
  {
    label: "今日新增处置",
    value: "5",
    delta: "+2",
    tone: "success",
    description: "包含新增封禁和新增解封记录",
    targetPath: "/unbans",
  },
];

export const latestViolations: ActivityItem[] = [
  {
    id: "lv-1",
    namespace: "prod-finance",
    summary: "CompliK 发现敏感关键字命中",
    time: "2026-04-06 19:42",
    tone: "danger",
    targetPath: "/namespaces/prod-finance",
  },
  {
    id: "lv-2",
    namespace: "ai-lab",
    summary: "Procscan 命中进程规则，需进一步确认",
    time: "2026-04-06 18:20",
    tone: "warn",
    targetPath: "/namespaces/ai-lab",
  },
  {
    id: "lv-3",
    namespace: "edge-demo",
    summary: "CompliK 检测到 URL 违规，已补充解释说明",
    time: "2026-04-06 17:10",
    tone: "info",
    targetPath: "/namespaces/edge-demo",
  },
];

export const latestActions: ActivityItem[] = [
  {
    id: "la-1",
    namespace: "prod-finance",
    summary: "新增封禁记录，操作人 Alice",
    time: "2026-04-06 19:55",
    tone: "warn",
    targetPath: "/bans",
  },
  {
    id: "la-2",
    namespace: "growth-ops",
    summary: "新增解封记录，操作人 Bob",
    time: "2026-04-06 16:18",
    tone: "success",
    targetPath: "/unbans",
  },
  {
    id: "la-3",
    namespace: "ai-lab",
    summary: "承诺书已更新，等待人工复核",
    time: "2026-04-06 13:06",
    tone: "info",
    targetPath: "/commitments",
  },
];

export const quickLinks: QuickLinkItem[] = [
  {
    title: "进入违规中心",
    description: "查看两类违规记录，并在右侧抽屉核对详情。",
    targetPath: "/violations",
  },
  {
    title: "查看封禁记录",
    description: "核对当前有效封禁，并补录新的封禁信息。",
    targetPath: "/bans",
  },
  {
    title: "维护项目配置",
    description: "查看配置类型、描述和 JSON 内容。",
    targetPath: "/configs",
  },
];

export const violations: ViolationRecord[] = [
  {
    id: "cv-1",
    apiId: 1,
    type: "complik",
    namespace: "prod-finance",
    detectorName: "keyword-detector",
    resourceName: "statement.pdf",
    host: "finance.internal",
    url: "https://finance.internal/statements/2026-q1",
    keywords: ["invoice", "confidential"],
    status: "open",
    detectedAt: "2026-04-06 19:42",
    description: "关键字规则命中，需人工确认是否违规。",
  },
  {
    id: "cv-2",
    apiId: 2,
    type: "complik",
    namespace: "edge-demo",
    detectorName: "url-detector",
    resourceName: "marketing-site",
    host: "edge.internal",
    url: "https://edge.internal/demo",
    keywords: ["download"],
    status: "open",
    detectedAt: "2026-04-06 17:10",
    description: "发现可疑 URL，待人工确认。",
  },
  {
    id: "pv-1",
    apiId: 3,
    type: "procscan",
    namespace: "ai-lab",
    processName: "nmap",
    podName: "scanner-7fd8f",
    nodeName: "node-cn-sh-02",
    message: "检测到高风险进程命中规则",
    status: "open",
    detectedAt: "2026-04-06 18:20",
    description: "进程规则命中，需要确认是否为测试行为。",
  },
  {
    id: "pv-2",
    apiId: 4,
    type: "procscan",
    namespace: "ops-tools",
    processName: "tcpdump",
    podName: "ops-agent-12aa",
    nodeName: "node-cn-bj-03",
    message: "抓包工具被标记为风险进程",
    status: "closed",
    detectedAt: "2026-04-06 11:32",
    description: "已确认是授权操作，记录保留用于追溯。",
  },
];

export const namespaceProfiles: NamespaceProfile[] = [
  {
    namespace: "prod-finance",
    violated: true,
    banned: true,
    commitmentUploaded: true,
    lastActionAt: "2026-04-06 19:55",
    commitment: {
      fileName: "prod-finance-commitment.pdf",
      fileUrl: "https://files.example.com/prod-finance-commitment.pdf",
      updatedAt: "2026-04-05 10:22",
    },
    recentViolations: violations.filter((item) => item.namespace === "prod-finance"),
    timeline: [
      {
        id: "t1",
        title: "新增封禁记录",
        description: "操作人 Alice，原因是敏感关键字命中。",
        time: "2026-04-06 19:55",
        tone: "warn",
      },
      {
        id: "t2",
        title: "出现新违规事件",
        description: "CompliK 关键字规则命中。",
        time: "2026-04-06 19:42",
        tone: "danger",
      },
      {
        id: "t3",
        title: "承诺书已上传",
        description: "当前承诺书版本已同步到记录中心。",
        time: "2026-04-05 10:22",
        tone: "success",
      },
    ],
  },
  {
    namespace: "ai-lab",
    violated: true,
    banned: false,
    commitmentUploaded: false,
    lastActionAt: "2026-04-06 18:20",
    recentViolations: violations.filter((item) => item.namespace === "ai-lab"),
    timeline: [
      {
        id: "t4",
        title: "出现新违规事件",
        description: "Procscan 进程规则命中，待进一步核查。",
        time: "2026-04-06 18:20",
        tone: "danger",
      },
    ],
  },
  {
    namespace: "growth-ops",
    violated: false,
    banned: false,
    commitmentUploaded: true,
    lastActionAt: "2026-04-06 16:18",
    commitment: {
      fileName: "growth-ops-commitment.pdf",
      fileUrl: "https://files.example.com/growth-ops-commitment.pdf",
      updatedAt: "2026-04-02 09:30",
    },
    recentViolations: [],
    timeline: [
      {
        id: "t5",
        title: "新增解封记录",
        description: "操作人 Bob，当前 namespace 已解除限制。",
        time: "2026-04-06 16:18",
        tone: "success",
      },
    ],
  },
];

export const configRecords: ConfigRecord[] = [
  {
    id: "cfg-1",
    configName: "project-config-demo",
    configType: "json",
    description: "默认项目配置示例",
    createdAt: "2026-04-06 10:42",
    updatedAt: "2026-04-06 10:42",
    value: JSON.stringify({ enabled: true, threshold: 3, reviewers: ["ops", "risk"] }, null, 2),
  },
  {
    id: "cfg-2",
    configName: "ban-policy",
    configType: "json",
    description: "封禁策略模板",
    createdAt: "2026-04-05 17:12",
    updatedAt: "2026-04-05 17:12",
    value: JSON.stringify({ duration_hours: 24, auto_notify: true }, null, 2),
  },
];

export const commitmentRecords: CommitmentRecord[] = [
  {
    id: "com-1",
    namespace: "prod-finance",
    fileName: "prod-finance-commitment.pdf",
    fileUrl: "https://files.example.com/prod-finance-commitment.pdf",
    createdAt: "2026-04-05 10:22",
    updatedAt: "2026-04-05 10:22",
  },
  {
    id: "com-2",
    namespace: "growth-ops",
    fileName: "growth-ops-commitment.pdf",
    fileUrl: "https://files.example.com/growth-ops-commitment.pdf",
    createdAt: "2026-04-02 09:30",
    updatedAt: "2026-04-02 09:30",
  },
];

export const banRecords: BanRecord[] = [
  {
    id: "ban-1",
    apiId: 1,
    namespace: "prod-finance",
    reason: "## 封禁说明\n- 敏感关键字命中\n- 需要临时封禁复核",
    screenshotUrls: ["https://files.example.com/prod-finance-ban-1.png"],
    operatorName: "Alice",
    banStartTime: "2026-04-06 19:55",
    banEndTime: "2026-04-07 19:55",
    createdAt: "2026-04-06 19:55",
    updatedAt: "2026-04-06 19:55",
    active: true,
  },
  {
    id: "ban-2",
    apiId: 2,
    namespace: "ops-tools",
    reason: "进程规则命中，待确认授权范围",
    screenshotUrls: [],
    operatorName: "Bob",
    banStartTime: "2026-04-05 09:10",
    banEndTime: "2026-04-05 18:10",
    createdAt: "2026-04-05 09:10",
    updatedAt: "2026-04-05 18:10",
    active: false,
  },
];

export const unbanRecords: UnbanRecord[] = [
  {
    id: "unban-1",
    apiId: 1,
    namespace: "growth-ops",
    operatorName: "Bob",
    createdAt: "2026-04-06 16:18",
    updatedAt: "2026-04-06 16:18",
  },
  {
    id: "unban-2",
    apiId: 2,
    namespace: "ops-tools",
    operatorName: "Alice",
    createdAt: "2026-04-05 18:30",
    updatedAt: "2026-04-05 18:30",
  },
];
