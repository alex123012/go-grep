import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

# .apply(lambda x: x.type(float))
df = pd.read_csv("result_table.csv", sep=', ', engine="python")
df.iloc[:, 1:] = df.iloc[:, 1:].apply(lambda x: x.astype(float))

plt.style.use("seaborn-poster")
print(df)
# columns: binary, time, mem, disk, files, pattern_len
fig, axes = plt.subplots(3, 2, figsize=(15, 20))
for column, ax in zip(df.columns[3:], axes):
    for ref, a in zip(df.columns[1:3], ax):
        style = "binary"
        if column == "pattern_len":
            style = "files"
        sns.lineplot(data=df, x=column, y=ref, hue="binary",
                     ax=a, markers=True, dashes=True, style=style)
        a.set_title(f"{ref}({column})")

fig.tight_layout()
fig.savefig("../.github/images/graphics.png", dpi=600)
