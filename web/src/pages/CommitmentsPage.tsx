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
  SurfaceCard,
} from "../components/ui";
import { useAppData } from "../contexts/AppDataContext";
import type { CommitmentRecord } from "../types";

export function CommitmentsPage() {
  const navigate = useNavigate();
  const { commitmentRecords, deleteCommitmentRecord } = useAppData();
  const [selected, setSelected] = useState<CommitmentRecord | null>(null);
  const [open, setOpen] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [pendingDelete, setPendingDelete] = useState<CommitmentRecord | null>(null);

  const rows = useMemo(() => {
    return commitmentRecords.filter((item) => item.namespace.toLowerCase().includes(keyword.toLowerCase()));
  }, [commitmentRecords, keyword]);

  return (
    <div className="page-container">
      <PageHeader
        kicker="Commitments"
        title="承诺书管理"
        description="按 namespace 查看承诺书记录，文件链接保持清晰可点，不在表格里铺长链接。"
        actions={<Button variant="primary" onClick={() => setOpen(true)}>新增承诺书记录</Button>}
      />

      <SurfaceCard>
        <div className="toolbar">
          <Field label="namespace 搜索">
            <Input placeholder="输入 namespace" value={keyword} onChange={(event) => setKeyword(event.target.value)} />
          </Field>
        </div>
      </SurfaceCard>

      <SurfaceCard className="data-table-wrap" padded={false}>
        {rows.length > 0 ? (
          <table className="data-table">
            <thead>
              <tr>
                <th>namespace</th>
                <th>文件名</th>
                <th>文件链接</th>
                <th>更新时间</th>
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
                      {item.fileName}
                    </button>
                  </td>
                  <td>
                    <a className="namespace-link" href={item.fileUrl} rel="noreferrer" target="_blank">
                      打开文件
                    </a>
                  </td>
                  <td>{item.updatedAt}</td>
                  <td>
                    <Button variant="ghost" onClick={() => setSelected(item)}>
                      查看
                    </Button>
                    <Button variant="danger" onClick={() => setPendingDelete(item)}>
                      删除
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <div style={{ padding: 20 }}>
            <EmptyState
              title="当前没有承诺书记录"
              description="可以直接新增一条承诺书记录。"
              action={<Button variant="primary" onClick={() => setOpen(true)}>新增承诺书记录</Button>}
            />
          </div>
        )}
      </SurfaceCard>

      <Drawer
        description="这里展示承诺书记录详情，并提供跳转到 namespace 详情的入口。"
        onClose={() => setSelected(null)}
        open={Boolean(selected)}
        title={selected ? selected.namespace : ""}
      >
        {selected ? (
          <>
            <DetailList
              items={[
                { label: "namespace", value: selected.namespace },
                { label: "文件名", value: selected.fileName },
                { label: "更新时间", value: selected.updatedAt },
                {
                  label: "文件链接",
                  value: (
                    <a className="namespace-link" href={selected.fileUrl} rel="noreferrer" target="_blank">
                      打开文件
                    </a>
                  ),
                },
              ]}
            />
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
        description="演示用新增承诺书表单。"
        onClose={() => setOpen(false)}
        open={open}
        title="新增承诺书记录"
      >
        <div className="panel-stack">
          <Field label="namespace">
            <Input placeholder="例如：prod-finance" />
          </Field>
          <Field label="文件名">
            <Input placeholder="例如：commitment.pdf" />
          </Field>
          <Field label="文件链接">
            <Input placeholder="https://files.example.com/commitment.pdf" />
          </Field>
          <div className="button-row">
            <Button variant="primary" onClick={() => setOpen(false)}>保存记录</Button>
            <Button variant="secondary" onClick={() => setOpen(false)}>取消</Button>
          </div>
        </div>
      </Modal>

      <ConfirmModal
        description={pendingDelete ? `删除后将从当前前端列表中移除 namespace ${pendingDelete.namespace} 的承诺书记录。` : ""}
        onClose={() => setPendingDelete(null)}
        onConfirm={() => {
          if (!pendingDelete) return;
          void deleteCommitmentRecord(pendingDelete.namespace).then(() => {
            if (selected?.id === pendingDelete.id) {
              setSelected(null);
            }
            setPendingDelete(null);
          });
        }}
        open={Boolean(pendingDelete)}
        title="删除承诺书记录"
      />
    </div>
  );
}
