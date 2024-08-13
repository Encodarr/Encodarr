import styles from "./MediaModal.module.scss";
import Modal from "../../modal/Modal";
import InputContainer from "../../inputs/inputContainer/InputContainer";

const MediaModal = ({ isOpen, setIsOpen, content, setContent, type }: any) => {
  const onClose = () => {
    setIsOpen(false);
  };

  const onSave = async () => {
    for (const key in content) {
      fetch(`/api/settings`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({ id: key, value: content[key] }),
      });
    }
    setIsOpen(false);
  };

  let title = "";
  if (type == "posters") title = "Poster Options";
  if (type == "table") title = "Table Options";
  if (type == "overview") title = "Overview Options";
  if (!isOpen) return null;
  return (
    <Modal
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      onSave={onSave}
      title={title}
      onClose={onClose}
    >
      <div className={styles.content}>
        <div className={styles.left}>
          {type === "posters" && (
            <>
              <InputContainer
                label="Poster Size"
                type="select"
                selected={content.media_poster_posterSize.value}
                onChange={(e: any) => {
                  setContent({
                    ...content,
                    media_poster_posterSize: e.target.value,
                  });
                }}
              >
                <option value={"small"}>Small</option>
                <option value={"medium"}>Medium</option>
                <option value={"large"}>Large</option>
              </InputContainer>
              <InputContainer
                label="Detailed Progress Bar"
                type="checkbox"
                checked={Boolean(
                  Number(content?.media_poster_detailedProgressBar.value)
                )}
                helpText="Show text on progress bar"
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_poster_detailedProgressBar: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Title"
                type="checkbox"
                checked={Boolean(content?.media_poster_showTitle.value)}
                helpText="Show series title under poster"
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_poster_showTitle: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Monitored"
                type="checkbox"
                checked={Boolean(content?.media_poster_showMonitored.value)}
                helpText="Show monitored status under poster"
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_poster_showMonitored: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Profile"
                type="checkbox"
                checked={Boolean(content?.media_poster_showProfile.value)}
                helpText="Show codec profile under poster"
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_poster_showProfile: e.target.checked,
                  })
                }
              />
            </>
          )}
          {type === "table" && (
            <>
              <InputContainer
                label="Show Type"
                type="checkbox"
                checked={Boolean(content?.media_table_showType.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showType: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Network"
                type="checkbox"
                checked={Boolean(content?.media_table_showNetwork.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showNetwork: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Year"
                type="checkbox"
                checked={Boolean(content?.media_table_showYear.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showYear: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Profile"
                type="checkbox"
                checked={Boolean(content?.media_table_showProfile.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showProfile: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Seasons"
                type="checkbox"
                checked={Boolean(content?.media_table_showSeasons.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showSeasons: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Episodes"
                type="checkbox"
                checked={Boolean(content?.media_table_showEpisodes.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showEpisodes: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Episode Count"
                type="checkbox"
                checked={Boolean(content?.media_table_showEpisodeCount.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showEpisodeCount: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Year"
                type="checkbox"
                checked={Boolean(content?.media_table_showYear.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showYear: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Size On Disk"
                type="checkbox"
                checked={Boolean(content?.media_table_showSizeOnDisk.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showSizeOnDisk: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Space Saved"
                type="checkbox"
                checked={Boolean(content?.media_table_showSizeSaved.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showSizeSaved: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Genre"
                type="checkbox"
                checked={Boolean(content?.media_table_showGenre.value)}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_table_showGenre: e.target.checked,
                  })
                }
              />
            </>
          )}
          {type == "overview" && (
            <>
              <InputContainer
                label="Poster Size"
                type="select"
                selected={content?.media_overview_posterSize}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_overview_posterSize: e.target.value,
                  })
                }
              >
                <option value={"small"}>Small</option>
                <option value={"medium"}>Medium</option>
                <option value={"large"}>Large</option>
              </InputContainer>
              <InputContainer
                label="Detailed Progress Bar"
                type="checkbox"
                checked={Boolean(
                  content?.media_overview_detailedProgressBar.value
                )}
                helpText={"Show text on progress bar"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_overview_detailedProgressBar: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Monitored"
                type="checkbox"
                checked={Boolean(content?.media_overview_showMonitored.value)}
                helpText={"Show monitored in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_overview_showMonitored: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Network"
                type="checkbox"
                checked={Boolean(content?.media_overview_showNetwork.value)}
                helpText={"Show network in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_overview_showNetwork: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Profile"
                type="checkbox"
                checked={Boolean(content?.media_overview_showProfile.value)}
                helpText={"Show profile in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_overview_showProfile: e.target.checked,
                  })
                }
              />

              <InputContainer
                label="Show Season Count"
                type="checkbox"
                checked={Boolean(content?.media_overview_showSeasonCount.value)}
                helpText={"Show season count in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_overview_showSeasonCount: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Path"
                type="checkbox"
                checked={Boolean(content?.media_overview_showPath.value)}
                helpText={"Show path in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_overview_showPath: e.target.checked,
                  })
                }
              />
              <InputContainer
                label="Show Size On Disk"
                type="checkbox"
                checked={Boolean(content?.media_overview_showSizeOnDisk.value)}
                helpText={"Show size on disk in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    media_overview_showSizeOnDisk: e.target.checked,
                  })
                }
              />
            </>
          )}
        </div>
      </div>
    </Modal>
  );
};
export default MediaModal;
