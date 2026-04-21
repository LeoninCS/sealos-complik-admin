import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Button,
  ConfirmModal,
  DetailList,
  Drawer,
  EmptyState,
  Field,
  Input,
  Modal,
  PageHeader,
  Select,
  SurfaceCard,
  StatusPill,
} from "../components/ui";
import { MarkdownInput } from "../components/MarkdownInput";
import { MarkdownRenderer } from "../components/MarkdownRenderer";
import { useAppData } from "../contexts/AppDataContext";
import { useManagedOperatorOptions } from "../hooks/useOperatorOptions";
import { buildBanScreenshotPreviewURL } from "../lib/api";
import { summarizeMarkdown } from "../lib/utils";
import type { BanRecord } from "../types";

export function BansPage() {
  const navigate = useNavigate();
  const { banRecords, configRecords, createBanRecord, deleteBanRecord, unbanRecords } = useAppData();
  const [selected, setSelected] = useState<BanRecord | null>(null);
  const [open, setOpen] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [operatorFilter, setOperatorFilter] = useState("");
  const [pendingDelete, setPendingDelete] = useState<BanRecord | null>(null);
  const [namespace, setNamespace] = useState("");
  const [reason, setReason] = useState("");
  const [banStartTime, setBanStartTime] = useState("");
  const [operatorName, setOperatorName] = useState("");
  const [screenshots, setScreenshots] = useState<File[]>([]);
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const { operatorConfigType, operatorOptions, operatorSource } = useManagedOperatorOptions(configRecords, [
    ...banRecords.map((item) => item.operatorName),
    ...unbanRecords.map((item) => item.operatorName),
  ]);

  const resetForm = () => {
    setNamespace("");
    setReason("");
    setBanStartTime("");
    setOperatorName("");
    setScreenshots([]);
  };

  const rows = useMemo(() => {
    return banRecords.filter((item) => {
      if (keyword && !item.namespace.toLowerCase().includes(keyword.toLowerCase())) {
        return false;
      }
      if (operatorFilter && item.operatorName !== operatorFilter) {
        return false;
      }
      return true;
    });
  }, [banRecords, keyword, operatorFilter]);

  const handleCreateBan = async () => {
    if (submitting) {
      return;
    }
    if (!namespace.trim() || !reason.trim() || !banStartTime.trim() || !operatorName.trim()) {
      setFormError("namespace、描述、开始时间、操作人均为必填。");
      return;
    }
    if (screenshots.length > 6) {
      setFormError("截图最多上传 6 张。");
      return;
    }

    setSubmitting(true);
    setFormError(null);
    try {
      await createBanRecord({
        namespace: namespace.trim(),
        reason: reason.trim(),
        banStartTime,
        operatorName: operatorName.trim(),
        screenshots,
      });
      setOpen(false);
      resetForm();
    } catch (err) {
      setFormError(err instanceof Error ? err.message : "新增封禁记录失败");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="page-container">
      <PageHeader
        kicker="Bans"
        title="封禁记录"
        description="录入和查看封禁记录，操作人从固定名单选择，描述支持 markdown，封禁截图支持一起上传。"
        actions={<Button variant="primary" onClick={() => setOpen(true)}>新增封禁</Button>}
      />

      <SurfaceCard>
        <div className="toolbar">
          <Field label="namespace">
            <Input placeholder="按 namespace 搜索" value={keyword} onChange={(event) => setKeyword(event.target.value)} />
          </Field>
          <Field label="操作人">
            <Select value={operatorFilter} onChange={(event) => setOperatorFilter(event.target.value)}>
              <option value="">全部操作人</option>
              {operatorOptions.map((option) => (
                <option key={option} value={option}>
                  {option}
                </option>
              ))}
            </Select>
          </Field>
        </div>
      </SurfaceCard>

      <SurfaceCard className="data-table-wrap" padded={false}>
        {rows.length > 0 ? (
          <table className="data-table">
            <thead>
              <tr>
                <th>namespace</th>
                <th>描述</th>
                <th>开始时间</th>
                <th>操作人</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => (
                <tr key={item.id}>
                  <td>
                    <button className="namespace-link table-row-button" onClick={() => navigate(`/namespaces/${item.namespace}`)} type="button">
                      {item.namespace}
                    </button>
                  </td>
                  <td>
                    <button className="table-row-button" onClick={() => setSelected(item)} type="button">
                      {summarizeMarkdown(item.reason, 72) || (item.screenshotUrls.length > 0 ? `截图 ${item.screenshotUrls.length} 张` : "-")}
                    </button>
                  </td>
                  <td>{item.banStartTime}</td>
                  <td>{item.operatorName}</td>
                  <td>
                    <StatusPill tone="warn">永久封禁</StatusPill>
                  </td>
                  <td>
                    <Button variant="ghost" onClick={() => setSelected(item)}>
                      查看
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <div style={{ padding: 20 }}>
            <EmptyState
              title="当前没有封禁记录"
              description="可以直接新增一条封禁记录。"
              action={<Button variant="primary" onClick={() => setOpen(true)}>新增封禁</Button>}
            />
          </div>
        )}
      </SurfaceCard>

      <Drawer
        description="这里展示 markdown 描述和截图附件。"
        onClose={() => setSelected(null)}
        open={Boolean(selected)}
        title={selected ? selected.namespace : ""}
      >
        {selected ? (
          <>
            <DetailList
              items={[
                { label: "namespace", value: selected.namespace },
                { label: "开始时间", value: selected.banStartTime },
                { label: "操作人", value: selected.operatorName },
                { label: "状态", value: "永久封禁" },
              ]}
            />
            <div className="ban-detail-section">
              <div className="detail-label">描述</div>
              <div className="detail-value">
                <MarkdownRenderer content={selected.reason} />
              </div>
            </div>
            <div className="ban-detail-section">
              <div className="detail-label">截图</div>
              {selected.screenshotUrls.length > 0 ? (
                <div className="screenshot-grid">
                  {selected.screenshotUrls.map((url, index) => {
                    const previewURL = buildBanScreenshotPreviewURL(url);
                    return (
                      <a className="screenshot-card" href={previewURL} key={`${url}-${index}`} rel="noreferrer" target="_blank">
                        <img alt={`封禁截图 ${index + 1}`} className="screenshot-image" loading="lazy" src={previewURL} />
                        <span className="screenshot-caption">截图 {index + 1}</span>
                      </a>
                    );
                  })}
                </div>
              ) : (
                <div className="muted-text">当前没有截图附件</div>
              )}
            </div>
            <div className="button-row" style={{ marginTop: 20 }}>
              <Button variant="secondary" onClick={() => navigate(`/namespaces/${selected.namespace}`)}>
                查看 namespace 详情
              </Button>
              <Button variant="danger" onClick={() => setPendingDelete(selected)}>
                删除记录
              </Button>
            </div>
          </>
        ) : null}
      </Drawer>

      <Modal
        description="描述使用 markdown 文本录入，截图会和封禁记录一起保存。"
        onClose={() => {
          setOpen(false);
          setFormError(null);
          resetForm();
        }}
        open={open}
        title="新增封禁"
      >
        <div className="panel-stack">
          <Field label="namespace">
            <Input placeholder="例如：prod-finance" value={namespace} onChange={(event) => setNamespace(event.target.value)} />
          </Field>
          <Field label="描述（Markdown）">
            <MarkdownInput
              placeholder={"例如：\n## 封禁说明\n- 违规链接已核实\n- 影响范围：prod-finance\n\n附上排查结论和后续动作。"}
              value={reason}
              onChange={setReason}
            />
          </Field>
          <Field label="开始时间">
            <Input type="datetime-local" value={banStartTime} onChange={(event) => setBanStartTime(event.target.value)} />
          </Field>
          <Field label="操作人">
            <Select value={operatorName} onChange={(event) => setOperatorName(event.target.value)}>
              <option value="">请选择操作人</option>
              {operatorOptions.map((option) => (
                <option key={option} value={option}>
                  {option}
                </option>
              ))}
            </Select>
          </Field>
          <div className="muted-text">
            {operatorSource === "config"
              ? `操作人名单来自配置类型 ${operatorConfigType}。`
              : `当前操作人名单来自历史记录。建议创建 config_type 为 ${operatorConfigType} 的配置，JSON 内容使用 {"operators":["张三","李四"]}。`}
          </div>
          <Field label="截图附件">
            <div className="upload-stack">
              <Input
                accept="image/png,image/jpeg,image/webp,image/gif"
                multiple
                type="file"
                onChange={(event) => setScreenshots(Array.from(event.target.files ?? []))}
              />
              <div className="muted-text">支持 PNG、JPG、WEBP、GIF，单次最多 6 张。</div>
              {screenshots.length > 0 ? (
                <div className="upload-list">
                  {screenshots.map((file) => (
                    <div className="upload-item" key={`${file.name}-${file.size}-${file.lastModified}`}>
                      <span>{file.name}</span>
                      <span className="muted-text">{Math.max(1, Math.round(file.size / 1024))} KB</span>
                    </div>
                  ))}
                </div>
              ) : null}
            </div>
          </Field>
          {formError ? <div className="muted-text" style={{ color: "#b42318" }}>{formError}</div> : null}
          <div className="button-row">
            <Button variant="primary" onClick={() => void handleCreateBan()}>
              {submitting ? "保存中..." : "保存封禁记录"}
            </Button>
            <Button
              variant="secondary"
              onClick={() => {
                setOpen(false);
                setFormError(null);
                resetForm();
              }}
            >
              取消
            </Button>
          </div>
        </div>
      </Modal>

      <ConfirmModal
        description={pendingDelete ? `删除后将从当前前端列表中移除 namespace ${pendingDelete.namespace} 的封禁记录。` : ""}
        onClose={() => setPendingDelete(null)}
        onConfirm={() => {
          if (!pendingDelete) return;
          void deleteBanRecord(pendingDelete.namespace).then(() => {
            if (selected?.id === pendingDelete.id) {
              setSelected(null);
            }
            setPendingDelete(null);
          });
        }}
        open={Boolean(pendingDelete)}
        title="删除封禁记录"
      />
    </div>
  );
}
