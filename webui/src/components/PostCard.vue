<script>
import ErrorMsg from './ErrorMsg.vue';
import Comment from './Comment.vue';

export default {
  name: 'PostCard',
  components: { ErrorMsg, Comment },
  props: {
    post: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      imageUrl: "",
      newCommentText: "",
      likesList: [],
      commentsList: [],
      showLikes: false,
      showComments: false,
      errorMessage: '',
      loading: false
    };
  },
  computed: {
    /**
     * Checks if the logged-in user is the owner of the current post.
     */
    isOwner() {
      return localStorage.getItem("Username") === this.post.username;
    }
  },
  methods: {
    /**
     * Build the direct URL for the post image stored in the backend.
     */
    buildImageUrl() {
      const baseUrl = this.$axios.defaults.baseURL;
      this.imageUrl = `${baseUrl}/users/${this.post.username}/posts/${this.post.postId}`;
    },

    /**
     * Toggles the Like status using PUT/DELETE requests.
     */
    async toggleLike() {
      this.loading = true;
      const token = localStorage.getItem("SessionToken");
      try {
        if (this.post.isLikedByMe) {
          await this.$axios.delete(`/users/${this.post.username}/posts/${this.post.postId}/likes`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.post.likesNumber--;
          this.likesList = this.likesList.filter(name => name !== localStorage.getItem("Username"));
          
        } else {
          await this.$axios.put(`/users/${this.post.username}/posts/${this.post.postId}/likes`, {}, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.post.likesNumber++;
          if (!this.likesList.includes(localStorage.getItem("Username"))) {
            this.likesList.push(localStorage.getItem("Username"));
          }
        }
        this.post.isLikedByMe = !this.post.isLikedByMe;
        this.loading = false;
      } catch (e) {
        this.loading = false;
        if (e.response) {
          switch (e.response.status) {
            case 401:
              localStorage.clear();
              this.$router.push("/login");
              break;
            case 403:
              this.errorMessage = "Azione non consentita.";
              break;
            case 404:
              this.errorMessage = "Risorsa non trovata.";
              break;
            default:
              this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
          }
        } else {
          // Handling network or connectivity issues
          this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
        }
      }
    },

    /**
     * Fetches the list of usernames who liked the post.
     */
    async getLikes() {
      this.loading = true;
      if (this.showLikes) { 
        this.showLikes = false; 
        this.loading = false;
        return; 
      }
      const token = localStorage.getItem("SessionToken");
      try {
        const response = await this.$axios.get(`/users/${this.post.username}/posts/${this.post.postId}/likes`, {headers: { Authorization: `Bearer ${token}` }});
        this.likesList = response.data;
        this.showLikes = true;
        this.showComments = false;
        this.loading = false;
        return;
      } catch (e) {
        this.loading = false;
        if (e.response) {
          switch (e.response.status) {
            case 401:
              localStorage.clear();
              this.$router.push("/login");
              break;
            case 403:
              this.errorMessage = "Azione non consentita.";
              break;
            case 404:
              this.errorMessage = "Risorsa non trovata.";
              break;
            default:
              this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
          }
        } else {
          // Handling network or connectivity issues
          this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
        }
        
      }
    },

    /**
     * Fetches comments for the post.
     */
    async getComments() {
      this.loading = true;
      if (this.showComments) { 
        this.showComments = false; 
        this.loading = false;
        return; 
      }
      const token = localStorage.getItem("SessionToken");
      try {
        const response = await this.$axios.get(`/users/${this.post.username}/posts/${this.post.postId}/comments`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.commentsList = response.data;
        this.showComments = true;
        this.showLikes = false;
        this.loading = false;
        return;
      } catch (e) {
        this.loading = false;
        if (e.response) {
          switch (e.response.status) {
            case 401:
              localStorage.clear();
              this.$router.push("/login");
              break;
            case 403:
              this.errorMessage = "Azione non consentita.";
              break;
            case 404:
              this.errorMessage = "Risorsa non trovata.";
              break;
            default:
              this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
          }
        } else {
          // Handling network or connectivity issues
          this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
        }
      }
    },

    /**
     * Sends a new comment to the backend and updates the local state.
     */
    async addComment() {
      if (!this.newCommentText.trim()) return;
      const token = localStorage.getItem("SessionToken");
      try {
        const response = await this.$axios.post(
          `/users/${this.post.username}/posts/${this.post.postId}/comments`,
          { content: this.newCommentText },
          { headers: { Authorization: `Bearer ${token}` } }
        );
        this.commentsList.push(response.data);
        this.post.commentsNumber++;
        this.newCommentText = "";
      } catch (e) {
        this.handleError(e, "Errore invio commento");
      }
    },

    /**
     * Deletes a specific comment and updates the local list.
     */
    async deleteComment(commentId) {
      this.loading = true;
      const token = localStorage.getItem("SessionToken");
      try {
        await this.$axios.delete(
          `/users/${this.post.username}/posts/${this.post.postId}/comments/${commentId}`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        this.commentsList = this.commentsList.filter(c => c.commentId !== commentId);
        this.post.commentsNumber--;
        this.loading = false;
      } catch (e) {
        this.loading = false;
        if (e.response) {
          switch (e.response.status) {
            case 401:
              localStorage.clear();
              this.$router.push("/login");
              break;
            case 403:
              this.errorMessage = "Azione non consentita.";
              break;
            case 404:
              this.errorMessage = "Risorsa non trovata.";
              break;
            default:
              this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
          }
        } else {
          // Handling network or connectivity issues
          this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
        }

      }
    },

    async deletePost() {
      if (!confirm("Sei sicuro di voler eliminare questo post?")) return;
      const token = localStorage.getItem("SessionToken");
      try {
        await this.$axios.delete(`/users/${this.post.username}/posts/${this.post.postId}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.$emit('post-deleted', this.post.postId);
      } catch (e) {
        this.handleError(e, "Errore eliminazione post");
      }
    }

    
  },
  mounted() {
    this.buildImageUrl();
  }
};
</script>

<template>
  <div class="post-card glass">
    <ErrorMsg v-if="errorMessage" :message="errorMessage" @close="errorMessage = ''" />

    <div class="post-header">
      <div class="user-info">
        <div class="avatar-mini">{{ post.username.charAt(0).toUpperCase() }}</div>
        <router-link :to="'/users/' + post.username" class="username-link">{{ post.username }}</router-link>
      </div>

      <button v-if="isOwner" @click="deletePost" class="btn-icon delete-btn"><i class="fa-solid fa-trash-can"></i></button>

    </div>

    <div class="post-media" @dblclick="toggleLike">
      <img v-if="imageUrl" :src="imageUrl" alt="Post photo" />

      <div v-else class="image-loader"><i class="fa-solid fa-circle-notch fa-spin"></i></div>
    </div>

    <div class="post-footer">
      <div class="actions-row">
        <div class="action-group">
          <button @click="toggleLike" class="btn-icon heart-icon" :class="{ 'liked': post.isLikedByMe }">

            <i :class="post.isLikedByMe ? 'fa-solid fa-heart' : 'fa-regular fa-heart'"></i>

          </button>
          <span class="stat-text" @click="getLikes">{{ post.likesNumber }}</span>
        </div>

        <div class="action-group">
          <button @click="getComments" class="btn-icon comment-icon"><i class="fa-regular fa-comment"></i></button>
          <span class="stat-text" @click="getComments">{{ post.commentsNumber }}</span>
        </div>
      </div>

      <div v-if="post.caption" class="post-caption">
        <span class="caption-username">{{ post.username }}</span> {{ post.caption }}
      </div>

      <div v-if="showLikes" class="expanded-list">
        <p class="list-title">Mi piace:</p>
        <div class="list-content">
          {{ likesList.length > 0 ? likesList.join(', ') : 'Ancora nessun mi piace.' }}
        </div>
      </div>

      <div v-if="showComments" class="expanded-list">
        <p class="list-title">Commenti</p>
        <div class="list-content scrollable">
          <Comment 
            v-for="c in commentsList" 
            :key="c.commentId" 
            :comment="c" 
            :postAuthor="post.username"
            @delete-comment="deleteComment"
          />
          <p v-if="commentsList.length === 0" class="empty-msg">Ancora nessun commento.</p>
        </div>

        <div class="add-comment-area">
          <input 
            type="text" v-model="newCommentText" 
            placeholder="Aggiungi un commento..." 
            @keyup.enter="addComment" class="comment-input"
          />
          <button class="btn-post-comment" :disabled="!newCommentText.trim()" @click="addComment">Invia</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.post-card {
  max-width: 500px; background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px); border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px; margin: 0 auto 30px; color: white;
  box-shadow: 0 10px 30px rgba(0,0,0,0.5); overflow: hidden;
}

.post-header { padding: 12px 15px; display: flex; justify-content: space-between; align-items: center; }
.user-info { display: flex; align-items: center; gap: 10px; }
.avatar-mini {
  width: 32px; height: 32px; border-radius: 50%;
  background: linear-gradient(45deg, #f09433, #dc2743, #bc1888);
  display: flex; align-items: center; justify-content: center; font-weight: bold; font-size: 0.8rem;
}
.username-link { color: white; text-decoration: none; font-weight: 600; font-size: 0.9rem; }

.post-media { width: 100%; aspect-ratio: 1/1; background: #000; display: flex; align-items: center; justify-content: center; }
.post-media img { width: 100%; height: 100%; object-fit: cover; }

.post-footer { padding: 12px 15px; }
.actions-row { display: flex; gap: 15px; margin-bottom: 10px; }
.action-group { display: flex; align-items: center; gap: 6px; }
.btn-icon { background: none; border: none; color: white; cursor: pointer; font-size: 1.3rem; transition: transform 0.1s; }
.btn-icon:hover { transform: scale(1.1); }
.liked i { color: #ff4757; }
.stat-text { font-size: 0.85rem; font-weight: 600; cursor: pointer; }
.delete-btn { color: rgba(255,255,255,0.4); font-size: 1rem; }

.post-caption { font-size: 0.9rem; margin-top: 5px; }
.caption-username { font-weight: 700; }

.expanded-list { margin-top: 15px; padding-top: 10px; border-top: 1px solid rgba(255,255,255,0.1); }
.list-title { font-weight: 700; color: #aaa; margin-bottom: 8px; font-size: 0.8rem; }
.scrollable { max-height: 180px; overflow-y: auto; padding-right: 5px; }

.add-comment-area {
  display: flex; align-items: center; background: rgba(255,255,255,0.05);
  border-radius: 8px; padding: 6px 12px; border: 1px solid rgba(255,255,255,0.1); margin-top: 10px;
}
.comment-input { flex: 1; background: transparent; border: none; color: white; font-size: 0.85rem; outline: none; }
.btn-post-comment { background: none; border: none; color: #0095f6; font-weight: 700; cursor: pointer; margin-left: 10px; }
.btn-post-comment:disabled { opacity: 0.3; }
.empty-msg { color: #666; font-style: italic; font-size: 0.8rem; }
</style>